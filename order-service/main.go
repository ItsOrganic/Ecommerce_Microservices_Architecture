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
	var err error
	err = db.Connect("mongodb://localhost:27017", "order-service", "orders")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	metrics.Init()
	utils.InitRedis()
	utils.InitMQ()
	defer utils.CloseMQ()

	router := gin.Default()
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/metrics", metrics.PrometheusHandler)
	router.GET("/orders", handler.GetOrders)
	router.POST("/order", handler.CreateOrder)
	router.GET("/order/:id", handler.GetOrder)
	router.PUT("/order/:id", handler.UpdateStatus)

	router.Run(":8083")
}
