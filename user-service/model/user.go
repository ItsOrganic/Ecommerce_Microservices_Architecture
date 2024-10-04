package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name     string             `json:"name" bson:"name" binding:"required"`
	Email    string             `json:"email" bson:"email" binding:"required,email"`
	Password string             `json:"password" bson:"password" binding:"required,min=6"`
}
