package utils

import (
	"log"

	"github.com/streadway/amqp"
)

var MQConnection *amqp.Connection
var MQChannel *amqp.Channel

// InitMQ initializes the RabbitMQ connection and channel
func InitMQ(uri string) error {
	var err error
	MQConnection, err = amqp.Dial(uri)
	if err != nil {
		return err
	}

	MQChannel, err = MQConnection.Channel()
	if err != nil {
		return err
	}

	log.Println("RabbitMQ connection and channel initialized")
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
	ch, err := MQConnection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
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
