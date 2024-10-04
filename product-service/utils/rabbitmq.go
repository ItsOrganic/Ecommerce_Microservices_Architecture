package utils

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func InitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
}
func CloseMQ() {
	conn.Close()
}

func EmitEvents(event string) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_ = ch.ExchangeDeclare(
		"user_service", // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err = ch.Publish(
		"user_service", // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(event),
		},
	); err != nil {
		panic(err)
	}
	log.Printf("sent event %s", event)
}
