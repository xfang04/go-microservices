package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

// Consumer is the type used for receiving AMPQ events
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer returns a new Consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup opens a channel and declares the exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload is the type used for pushing events to RabbitMQ
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen will listen for all new queue publications
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Println(err)
		}
	}(ch)

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			getExchangeName(),
			false,
			nil,
		)

		if err != nil {
			log.Println(err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			// get the JSON payload and unmarshal it into a variable
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// do something with the payload
			go handlePayload(payload)
		}
	}()

	log.Printf("[*] Waiting for message [Exchange, Queue][%s, %s].", getExchangeName(), q.Name)
	<-forever
	return nil
}

// handlePayload takes an action based on the name of an event in the queue
func handlePayload(payload Payload) {
	// logic to process payload goes in here
	switch payload.Name {
	case "auth", "authentication":
		// we are trying to authenticate someone
		err := authenticate(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Got payload of", payload)
	}
}

// authenticate is a stub that we'll never actually use, but it is here
// as we get used to how to interact with services
func authenticate(payload Payload) error {
	// TODO actually authenticate via JSON
	log.Printf("Got payload of %v", payload)
	return nil
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
