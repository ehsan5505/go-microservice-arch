package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

type Consumer struct {
	conn 			*amqp.Connection
	queueName	string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error ) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{},err
	}
	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// for push messages in the Rabbit
type Payload struct {
	Name 		string 	`json:"name"`
	data		string	`json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	// we have the channel, and close it when complete the use
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _,s := range topics {
		ch.QueueBind(
			q.Name,						// name,
			s,								// topic / key
			"logs_topic",			// exchange string
			false,						// no wait
			nil,							// args
		)
		if err != nil {
			return err
		}
	}

	// https://pkg.go.dev/github.com/rabbitmq/amqp091-go#Channel.Consume
	messages, err := ch.Consume(
		q.Name,			// queue
		"",					// consumer
		true,				// autoAcknowledge
		false,			// exclusive?
		false,			// noLocal?
		false,			// no wait?
		nil					// args
	);
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d:= range messages {
			var payload Payload 
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for Message [Exchange, Queue] [logs_topic, %s]\n",q.Name)<-forever

	return nil

}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log","event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate
	// we have to multiple cases you want for the listener to listen and trigger the event
	default:
		err:= logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST",logServiceURL,bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type","application/json")

	client := &http.Client()

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCode {
		return err
	}

	return nil 

}