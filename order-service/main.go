package main

import (
	"log"
	"order-service/db"
	"order-service/handler"
	"order-service/metrics"
	"order-service/middleware"
	"order-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start the server
	err := db.Connect("mongodb://mongodb:27017/order-service", "order-service", "orders")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	metrics.Init()

	utils.InitMQ()
	defer utils.CloseMQ()

	router := gin.Default()
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/metrics", metrics.PrometheusHandler)
	router.POST("/order", handler.CreateOrder)
	router.GET("/order/:id", handler.GetOrder)
	router.PUT("/order/:id", handler.UpdateStatus)

	router.Run(":8083")
}
