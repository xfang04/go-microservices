package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/lib/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {

		}
	}(rabbitConn)

	// start listening for messages

	// create a new consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// consumer.Listen watches the queue and consumes events for all the provided topics.
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

// connect tries to connect to RabbitMQ, and delays between attempts.
// If we can't connect after 5 tries (with increasing delays), return an error
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	var rabbitURL = "amqp://user:pass@rabbitmq"
	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			fmt.Println("RabbitMQ is ready!")
			connection = c
			break
		}

		if counts > 5 {
			// if we can't connect after five tries, something is wrong...
			fmt.Println(err)
			return nil, err
		}
		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
