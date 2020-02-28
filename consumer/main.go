package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"

	"github.com/kravcs/go_rabbit"
)

func main() {
	// set up connection to RabbitMQ
	connection, err := amqp.Dial(go_rabbit.Config.AMQPConnectionURL)
	go_rabbit.HandleError(err, "Can't connect to AMPQ")
	defer connection.Close()

	// set up a RabbitMQ channel to communicate
	amqpChannel, err := connection.Channel()
	go_rabbit.HandleError(err, "Can't connect to AMPQ Channel")
	defer amqpChannel.Close()

	// declare our `add` queue to get messages from
	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	go_rabbit.HandleError(err, "Could not declare `add` queue")

	// ask RabbitMQ to deliver new messages only when the worker has acknowledged previous message
	err = amqpChannel.Qos(1, 0, false)
	go_rabbit.HandleError(err, "Could not configure QoS")

	// set up consumer of messages from our queue
	messageChannel, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)
	go_rabbit.HandleError(err, "Could not configure QoS")

	forever := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Recieved a message: %s", d.Body)

			addTask := go_rabbit.AddTask{}
			err := json.Unmarshal(d.Body, &addTask)
			go_rabbit.HandleError(err, "Error decoding JSON")

			log.Printf("Result of %d + %d is : %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			err = d.Ack(false)
			go_rabbit.HandleError(err, "Error acknowledging message")
			if err == nil {
				log.Printf("Acknowledged a message")
			}
		}
	}()

	// Stop for programm termination
	<-forever
}
