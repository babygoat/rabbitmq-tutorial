package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	s := ""
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	return s
}

// Create a durable queue

func main() {
	// version-validation and authentication to the RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// setup the real connection channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// the queue will only be created if it does not exist already
	q, err := ch.QueueDeclare(
		"task_queue", // queue name
		true,         // durable (mesage will not be restored when the server restarts)
		false,        // delete when unused (no consumer)
		false,        // exclusive queue (only accessible to the declarer)
		false,        // no-wait (?)
		nil,          // arguments
	)

	failOnError(err, "Failed to declare a queue")

	body := bodyFrom(os.Args)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory (might not be published if no queue is bound to the exchange)
		false,  // immediate (might not be published if no consumer is on the matched queue)
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
