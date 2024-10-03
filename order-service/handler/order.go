package handler

import (
	"context"
	"fmt"
	"net/http"
	"order-service/db"
	"order-service/model"
	"order-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateOrder(c *gin.Context) {
	// Create a new order
	var order model.Order
	err := c.BindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.Status = "pending"
	order.CreatedAt = time.Now().Unix()
	order.ID = primitive.NewObjectID()

	// Insert the order to the database
	_, err = db.MI.Collection.InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Emit an event
	event := fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, order.ID, order.Status)
	utils.EmitEvents(event)

	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	var order model.Order

	// Extract the order ID from the URL parameter
	orderId := c.Param("id")

	// Convert the order ID from string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Create a filter to find the order by ID
	filter := bson.M{"_id": objectId}

	// Fetch the order from the collection
	err = db.MI.Collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving order"})
		return
	}

	// Return the order if found
	c.JSON(http.StatusOK, order)
}

func UpdateStatus(c *gin.Context) {
	userId := c.Param("id")
	var order model.Order
	err := c.BindJSON(&order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the order status
	if order.Status == "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Already Shipped"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filter := bson.M{"_id": objectId}
	status := bson.M{"$set": bson.M{"status": order.Status}}
	_, err = db.MI.Collection.UpdateOne(context.TODO(), filter, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating order status"})
		return
	}
	utils.EmitEvents(fmt.Sprintf(`{"order_id": "%s", "status": "%s"}`, userId, order.Status))
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
