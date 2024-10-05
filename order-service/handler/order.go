package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/db"
	"order-service/model"
	"order-service/utils" // Ensure this path is correct relative to your project structure
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateOrder handles the creation of a new order
func CreateOrder(c *gin.Context) {
	var order model.Order
	// Bind JSON data to order struct
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Check the product inventory via HTTP request to the product service
	productResp, err := http.Get(fmt.Sprintf("http://localhost:8082/product/%s", order.ProductName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error sending request to product service: %v", err)})
		return
	}
	defer productResp.Body.Close()
	if productResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Decode the product details
	var product model.Product
	if err := json.NewDecoder(productResp.Body).Decode(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error decoding product response: %v", err)})
		return
	}

	// Check if enough inventory is available
	if product.Quantity < order.Quantity {
		c.JSON(http.StatusConflict, gin.H{"error": "insufficient inventory"})
		return
	}

	// Set additional order fields
	order.CreatedAt = time.Now().Format(time.RFC3339)
	order.Price = product.Price // Use the product's price

	// Save the new order to your MongoDB database
	insertResult, err := db.MI.DB.Collection("orders").InsertOne(context.TODO(), order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error creating order: %v", err)})
		return
	}
	log.Printf("Inserted a single document: %v", insertResult.InsertedID)

	// Emit an event to RabbitMQ
	event := map[string]interface{}{
		"product_name": order.ProductName,
		"quantity":     order.Quantity,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error marshaling event: %v", eventJSON)})
		return
	}
	utils.EmitEvents("Order Created")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error emitting event: %v", err)})
		return
	}

	// Update the product inventory
	err = utils.UpdateProductInventory(order.ProductName, -order.Quantity)
	if err != nil {
		log.Printf("Error updating product inventory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error updating product inventory: %v", err)})
		return
	}

	utils.EmitEvents("Order_created")

	// Return the created order response
	c.JSON(http.StatusCreated, order)
}

// GetOrder retrieves an order by product name
func GetOrder(c *gin.Context) {
	productName := c.Param("productName")
	var order model.Order

	// Check Redis cache first
	val, err := utils.RDB.Get(context.Background(), "order:"+productName).Result()
	if err == nil {
		// Cache hit
		log.Println("Cache hit")
		if err := json.Unmarshal([]byte(val), &order); err == nil {
			c.JSON(http.StatusOK, order)
			return
		}
	}

	// Cache miss, query the database
	err = db.MI.DB.Collection("orders").FindOne(context.TODO(), bson.M{"productName": productName}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching order"})
		return
	}

	// Store result in Redis cache
	data, err := json.Marshal(order)
	if err == nil {
		err = utils.RDB.Set(context.Background(), "order:"+productName, data, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Error setting cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, order)
}

// UpdateStatus updates the status of an order by product name
func UpdateStatus(c *gin.Context) {
	productName := c.Param("name")
	var order model.Order

	// Bind JSON data to order struct
	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the order status can be updated
	if order.Status == "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order Already Shipped"})
		return
	}

	// Create a filter to find the order by product name
	filter := bson.M{"productName": productName}
	update := bson.M{"$set": bson.M{"status": order.Status}}

	// Update the order status
	_, err := db.MI.DB.Collection("orders").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating order status"})
		return
	}

	utils.EmitEvents("order_exchange")

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func GetOrders(c *gin.Context) {
	var orders []model.Order

	// Query the database
	cursor, err := db.MI.DB.Collection("orders").Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching orders"})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var order model.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}

	c.JSON(200, orders)
}
