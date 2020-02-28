package go_rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

type Configuration struct {
	AMQPConnectionURL string
}

type AddTask struct {
	Number1 int
	Number2 int
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@localhost:5672",
}

var AmqpConnection *amqp.Connection
var AmqpChannel *amqp.Channel

func init() {
	// set up connection to RabbitMQ
	AmqpConnection, err := amqp.Dial(Config.AMQPConnectionURL)
	HandleError(err, "Can't connect to AMPQ")
	defer AmqpConnection.Close()

	// set up a RabbitMQ channel to communicate
	AmqpChannel, err := AmqpConnection.Channel()
	HandleError(err, "Can't connect to AMPQ Channel")
	defer AmqpChannel.Close()

	// create an exchange that will bind to the queue to send and receive messages
	err = AmqpChannel.ExchangeDeclare("calc", "direct", true, false, false, false, nil)
	HandleError(err, "Could not declare `calc` queue")

	// declare our `calculate` queue to send messages to it
	_, err = AmqpChannel.QueueDeclare("calculate", true, false, false, false, nil)
	HandleError(err, "Could not declare `calculate` queue")

	// bind exchange to queue
	err = AmqpChannel.QueueBind("calculate", "#", "calc", false, nil)
	HandleError(err, "Could not bind `calculate` queue to `calc` exchange")
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
