package model

type Order struct {
	ProductName string  `json:"name" bson:"name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	Price       float64 `json:"price" bson:"price"`
	Status      string  `json:"status" bson:"status"`
	CreatedAt   string  `json:"created_at" bson:"created_at"`
}
