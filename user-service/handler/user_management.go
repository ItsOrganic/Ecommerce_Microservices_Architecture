package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"user-service/db"
	"user-service/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UpdateProfile allows a user to update their name and password
func UpdateProfile(c *gin.Context) {
	var update UpdateProfileRequest
	userId := c.Param("id") // Ensure you retrieve the user ID correctly from the URL parameter

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert userId from string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Prepare the password if it's being updated
	if update.Password != "" {
		hashedPassword, err := hashPassword(update.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		update.Password = hashedPassword // Use the hashed password for the update
	}

	// Build the update query
	updateFields := bson.M{}
	if update.Name != "" {
		updateFields["name"] = update.Name
	}
	if update.Password != "" {
		updateFields["password"] = update.Password
	}

	// Perform the update
	_, err = db.MI.DB.Collection("users").UpdateOne(context.TODO(),
		bson.M{"_id": objectId},
		bson.M{"$set": updateFields})
	if err != nil {
		log.Println("Error updating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	utils.EmitEvents(fmt.Sprintf("user updated %s", update.Name))

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetProfile(c *gin.Context) {
	userId := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user bson.M
	err = db.MI.DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// hashPassword hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
