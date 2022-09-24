package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// https://pkg.go.dev/github.com/rabbitmq/amqp091-go#Channel.ExchangeDeclare
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic",     // name
		"topic",					// kind
		true,							//durable?
		false,						// autodelete?
		false,						// internal?
		false,						// no-wait?
		nil,								// args?
	)
}

// https://pkg.go.dev/github.com/rabbitmq/amqp091-go#Channel.QueueDeclare
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	
	return ch.QueueDeclare(
		"",					// name
		false,			// durable?
		false,			//autoDelete?
		true,				//exclusive?
		false,			//no wait 
		nil,					// args

	) 
}