package main

import (
	"log"
	"order-service/db"
	"order-service/handler"
	"order-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start the server
	err := db.Connect("mongodb://localhost:27017", "order-service", "orders")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	utils.InitMQ()
	defer utils.CloseMQ()

	router := gin.Default()
	router.POST("/order", handler.CreateOrder)
	router.GET("/order/:id", handler.GetOrder)
	router.PUT("/order/:id", handler.UpdateStatus)

	router.Run(":8083")
}
