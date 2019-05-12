package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// version-validation and authentication to the RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// setup the real connection channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // Exchange name
		"fanout", // Exchange type
		true,     // Durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // auto-generated queue name
		false, // non-durable
		false, // delete when unused (no consumer)
		true,  // exclusive queue (only accessible to the declarer)
		false, // no-wait (?)
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",
		"logs", // exchange name
		false,
		nil,
	)

	failOnError(err, "Failed to bind a queue with an exchange")

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Recevied a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
