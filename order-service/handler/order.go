package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"order-service/db"
	"order-service/model"
	"order-service/utils"
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

	// Step 2: Decode the product details
	var product model.Product
	if err := json.NewDecoder(productResp.Body).Decode(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error decoding product response: %v", err)})
		return
	}

	// Step 3: Check if enough inventory is available
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

	// Step 4: Update the inventory in the product service
	updatedQuantity := product.Quantity - order.Quantity // Calculate updated quantity
	updatedQuantityURL := fmt.Sprintf("http://localhost:8082/product/%s", order.ProductName)

	// Create the JSON payload for the inventory update
	updatedQuantityJSON, err := json.Marshal(map[string]interface{}{
		"product_name": order.ProductName,
		"quantity":     updatedQuantity,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error marshaling updated quantity: %v", err)})
		return
	}
	log.Print(updatedQuantityJSON)

	// Create a new request to update the inventory
	req, err := http.NewRequest(http.MethodPut, updatedQuantityURL, bytes.NewBuffer(updatedQuantityJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error creating request to update inventory: %v", err)})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error updating inventory: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Log the status code for debugging
	fmt.Printf("Inventory Update Response Status: %s\n", resp.Status)

	if resp.StatusCode != http.StatusOK {
		// Read the response body for more details
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error updating inventory: %s\n", body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update inventory"})
		return
	}

	utils.EmitEvents("order created")

	// Return the created order response
	c.JSON(http.StatusCreated, order)
}

// GetOrder retrieves an order by product name
func GetOrder(c *gin.Context) {
	var order model.Order

	// Find the order by product name
	cursor, err := db.MI.DB.Collection("orders").Find(context.TODO(), bson.M{"productName": order.ProductName})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching order"})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		cursor.Decode(&order)
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

	utils.EmitEvents("order status updated")

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func GetOrders(c *gin.Context) {
	var orders []model.Order

	// Check Redis cache first
	val, err := utils.RDB.Get(context.Background(), "orders").Result()
	if err == nil {
		// Cache hit
		log.Println("Cache hit")
		if err := json.Unmarshal([]byte(val), &orders); err == nil {
			c.JSON(200, orders)
			return
		}
	}

	// Cache miss, query the database
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

	// Store result in Redis cache
	data, err := json.Marshal(orders)
	if err == nil {
		err = utils.RDB.Set(context.Background(), "orders", data, 0).Err()
		if err != nil {
			log.Printf("Error setting cache: %v", err)
		}
	}

	c.JSON(200, orders)
}
