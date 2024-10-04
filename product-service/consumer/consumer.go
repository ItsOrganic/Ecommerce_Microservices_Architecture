package consumer

import (
	"encoding/json"
	"log"
	"product-service/db"
	"product-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

// OrderEvent represents the structure of the order event
type OrderEvent struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

// StartConsumer starts the RabbitMQ consumer
func StartConsumer() {
	msgs, err := utils.MQChannel.Consume(
		"order_queue", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var event OrderEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding event: %v", err)
				continue
			}
			updateInventory(event)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}

// updateInventory updates the inventory based on the order event
func updateInventory(event OrderEvent) {
	filter := bson.M{"name": event.ProductName}
	update := bson.M{"$inc": bson.M{"quantity": -event.Quantity}}

	_, err := db.MI.DB.Collection("products").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Product not found: %s", event.ProductName)
			return
		}
		log.Printf("Error updating inventory: %v", err)
		return
	}

	log.Printf("Inventory updated for product: %s", event.ProductName)
}
