package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	ProductID string             `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
}
