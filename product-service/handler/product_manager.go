package handler

import (
	"context"
	"product-service/db"
	"product-service/model"
	"product-service/utils"

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

	utils.EmitEvents("product_created")
	c.JSON(200, gin.H{"message": "Product created successfully", "data": product})
}

func UpdateInventory(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filter := bson.D{{Key: "name", Value: product.ProductName}}
	update := bson.M{
		"$set": bson.M{
			"quantity": product.Quantity,
		},
	}
	_, err := db.MI.DB.Collection("products").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating product"})
		return
	}

	utils.EmitEvents("inventory_updated")

	c.JSON(200, gin.H{"message": "Inventory updated successfully"})
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
	utils.EmitEvents("product_deleted")

	c.JSON(200, gin.H{"message": "Product deleted successfully"})
}

func GetProducts(c *gin.Context) {
	var products []model.Product
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
	filter := bson.D{{Key: "name", Value: productName}}
	var product model.Product
	err := db.MI.DB.Collection("products").FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(200, product)
}
