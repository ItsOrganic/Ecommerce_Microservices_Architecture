package utils

import (
	"log"

	"github.com/streadway/amqp"
)

var MQConnection *amqp.Connection
var MQChannel *amqp.Channel

// InitMQ initializes the RabbitMQ connection and channel
func InitMQ() error {
	var err error
	MQConnection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	MQChannel, err = MQConnection.Channel()
	if err != nil {
		return err
	}
	log.Println("RabbitMQ connection and channel initialized")

	// Declare the exchange
	err = MQChannel.ExchangeDeclare(
		"order_exchange", // name
		"fanout",         // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	return nil
}

// CloseMQ closes the RabbitMQ connection and channel
func CloseMQ() {
	if MQChannel != nil {
		MQChannel.Close()
	}
	if MQConnection != nil {
		MQConnection.Close()
	}
	log.Println("RabbitMQ connection and channel closed")
}

// EmitEvent publishes an event to the specified RabbitMQ exchange
func EmitEvent(exchangeName, event string) error {
	err := MQChannel.Publish(
		exchangeName, // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(event),
		},
	)
	return err
}
