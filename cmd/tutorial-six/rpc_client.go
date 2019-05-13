package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) int {
	s := ""
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}

	n, err := strconv.Atoi(s)
	failOnError(err, "Failed to convert arg to integer")
	return n
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func fibonacciRPC(n int) (res int, err error) {
	// version-validation and authentication to the RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// setup the real connection channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false, // delete when unused
		true,  // exclusivve
		false,
		nil, // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",   // auto-generated consumer name (ConsumerTag)
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		},
	)
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	n := bodyFrom(os.Args)

	log.Printf(" [x] Requesting fib(%d)", n)
	res, err := fibonacciRPC(n)
	failOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Get %d", res)
}
