package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"product-service/db"
	"product-service/model"
	"product-service/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}}, // 1 for ascending order
		Options: options.Index().SetUnique(true),
	}
	db.MI.DB.Collection("products").Indexes().CreateOne(context.TODO(), indexModel)

	_, err := db.MI.DB.Collection("products").InsertOne(c, product)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating product/ Product already exists"})
		return
	}

	utils.EmitEvents("Product Created")

	c.JSON(200, gin.H{"message": "Product created successfully", "data": product})
}

func UpdateProduct(c *gin.Context) {
	productName := c.Param("name")
	var updateData struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"name": productName}
	update := bson.M{"$inc": bson.M{"quantity": updateData.Quantity}}

	result, err := db.MI.DB.Collection("products").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	utils.EmitEvents("product_updated")

	c.JSON(http.StatusOK, gin.H{"message": "product inventory updated"})
}

func DeleteProduct(c *gin.Context) {
	productName := c.Param("name")
	filter := bson.D{{Key: "name", Value: productName}}

	find, err := db.MI.DB.Collection("products").Find(context.TODO(), filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error finding product"})
		return
	}
	if !find.Next(context.Background()) {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	_, err = db.MI.DB.Collection("products").DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting product"})
		return
	}
	utils.EmitEvents("Product Deleted")

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}

func GetProducts(c *gin.Context) {
	var products []model.Product

	// Query the database
	cursor, err := db.MI.DB.Collection("products").Find(context.Background(), bson.D{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Error fetching products"})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var product model.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	c.JSON(200, products)
}

func GetProduct(c *gin.Context) {
	productName := c.Param("name")
	var product model.Product

	// Check Redis cache first
	val, err := utils.RDB.Get(context.Background(), "product:"+productName).Result()
	if err == nil {
		// Cache hit
		log.Println("Cache hit")
		if err := json.Unmarshal([]byte(val), &product); err == nil {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	// Cache miss, query the database
	err = db.MI.DB.Collection("products").FindOne(context.TODO(), bson.M{"name": productName}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching product"})
		return
	}

	// Store result in Redis cache
	data, err := json.Marshal(product)
	if err == nil {
		err = utils.RDB.Set(context.Background(), "product:"+productName, data, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Error setting cache: %v", err)
		}
	}

	c.JSON(http.StatusOK, product)
}
