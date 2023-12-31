package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	rabbitConn, err := connectToRabbit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(rabbitConn *amqp.Connection) {
		err := rabbitConn.Close()
		if err != nil {

		}
	}(rabbitConn)
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connectToRabbit tries to connect to RabbitMQ, for up to 30 seconds
func connectToRabbit() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	var rabbitURL = "amqp://user:pass@rabbitmq"
	// don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println(err)
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
