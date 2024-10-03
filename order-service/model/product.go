package model

type Product struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	ProductName string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
	Quantity    int     `json:"quantity" bson:"quantity"`
}
