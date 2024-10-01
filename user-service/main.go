package main

import (
	"log"
	"user-service/db"
	"user-service/handler"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func main() {
	var err error
	err = db.Connect("mongodb://localhost:27017", "user-service", "users")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	utils.InitMQ()
	defer utils.CloseMQ()

	router := gin.Default()
	router.POST("/register", handler.RegisterUser)
	router.POST("/login", handler.AuthenticateUser)
	router.PUT("/profile/:id", handler.UpdateProfile)
	router.Run(":8080")

}
