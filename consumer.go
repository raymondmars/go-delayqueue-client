package godelayqueueclient

import (
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProcessor func(msgs <-chan amqp.Delivery)

func ReceiveRabbitMQMessage(queueName string, processor RabbitMQProcessor) error {
	configSupplier := getEvnWithDefaultVal("QUEUE_SUPPLIER", "")
	if configSupplier == "rabbitmq" {
		amqpURI := getEvnWithDefaultVal("AMQP_URI", "amqp://guest:guest@localhost:5672/")
		// create connection
		connection, err := amqp.Dial(amqpURI)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer connection.Close()

		// create channel
		channel, err := connection.Channel()
		failOnError(err, "Failed to open a channel")
		defer channel.Close()

		// create queue
		q, err := channel.QueueDeclare(
			string(queueName), // name
			true,              // durable
			false,             // delete when unused
			false,             // exclusive
			false,             // no-wait
			nil,               // arguments
		)
		failOnError(err, "Failed to declare a queue")
		// subscribe message
		msgs, err := channel.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		//will wait here
		processor(msgs)

		return nil
	} else {
		return errors.New("QUEUE_SUPPLIER is not rabbitmq")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
