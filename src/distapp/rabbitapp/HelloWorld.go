package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go server()
	go client()

	var a string
	fmt.Scanln(&a)
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}
	for {
		ch.Publish("", q.Name, false, false, msg)
	}
}

func client() {
	conn, ch, q := getQueue()

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	failOnError(err, "failed to recieve the message")

	for msg := range msgs {
		log.Printf("Recieved messages with message: %s", msg.Body)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "failed connecting to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "failed opening a channel")
	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "failed to declare a queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
