package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
	"fmt"
	"log"
	"os"
	"time"
	"math"
)

func main(){

	// create the connection with RabbitMQ Service
	rabbitConn, err := connect()
	if err != nil {
		log.Println("Error to connect")
		os.Exit(1)
	}
	defer rabbitConn.Close() // close the connection for now
	
	// Listen to the Services
	log.Println("RabbitMQ Listening for the event request")

	// Consume the data
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Watch consume and trigger event
	err = consumer.Listen([]string{"log.INFO","log.WARNING","log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection,error) {
	// define the contstant
	var waitTime = 1 * time.Second
	var conn *amqp.Connection
	var count int64

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Rabbit MQ Not Ready...")
			count++
		}else {
			log.Println("Connected Successfully to RabbitMQ")
			conn = c
			break
		}

		if (count > 5 ) {
			log.Println(err)
			return  nil, err
		}

		waitTime = time.Duration(math.Pow(float64(waitTime),2))
		log.Println("Retry for the Connection")
		time.Sleep(waitTime)
	}

	// atlast connection was stablished
	return conn,nil

}
