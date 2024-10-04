package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"user-service/db"
	"user-service/model"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(c *gin.Context) {
	// validate user credentials
	// generate JWT token
	// return token
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	// find user by email
	err := db.MI.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	//Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	//Generate JWT token
	token, err := utils.GenerateToken(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	userJson, _ := json.Marshal(user)
	err = utils.EmitEvent("user authenticated", string(userJson))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error emitting event"})
		return
	}

	utils.EmitEvent("User authenticated", string(userJson))
	c.JSON(http.StatusOK, gin.H{"token": token, "message": "User authenticated"})

}
