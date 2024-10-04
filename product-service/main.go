package main

import (
	"log"
	"product-service/db"
	"product-service/handler"
	"product-service/metrics"
	"product-service/middleware"
	"product-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error
	err = db.Connect("mongodb://localhost:27017", "product-service", "products")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	metrics.Init()
	utils.InitRedis()
	utils.InitMQ()
	defer utils.CloseMQ()

	router := gin.Default()
	router.Use(middleware.PrometheusMiddleware())
	router.GET("/metrics", metrics.PrometheusHandler)
	router.POST("/product", handler.CreateProduct)
	router.GET("/product/:name", handler.GetProduct)
	router.GET("/products", handler.GetProducts)
	router.PUT("/product/:name", handler.UpdateProduct)
	router.DELETE("/product/:name", handler.DeleteProduct)
	router.Run(":8082")
}
