package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"

	"github.com/kravcs/go_rabbit"
)

func init() {
    // set up connection to RabbitMQ
	conn, err := amqp.Dial(go_rabbit.Config.AMQPConnectionURL)
	go_rabbit.HandleError(err, "Can't connect to AMPQ")
	defer conn.Close()

	// set up a RabbitMQ channel to communicate
	amqpChannel, err := conn.Channel()
	go_rabbit.HandleError(err, "Can't connect to AMPQ Channel")
	defer amqpChannel.Close()

	// create an exchange that will bind to the queue to send and receive messages
	err = amqpChannel.ExchangeDeclare("calc", "direct", true, false, false, false, nil)
	go_rabbit.HandleError(err, "Could not declare `calc` queue")

	// declare our `calculate` queue to send messages to it
	_, err = amqpChannel.QueueDeclare("calculate", true, false, false, false, nil)
	go_rabbit.HandleError(err, "Could not declare `calculate` queue")

	// bind exchange to queue
	err = amqpChannel.QueueBind("calculate", "#", "calc", false, nil)
	go_rabbit.HandleError(err, "Could not bind `calculate` queue to `calc` exchange")
}

func main() {
	
	// prepare message body to be send
	rand.Seed(time.Now().UnixNano())
	addTask := go_rabbit.AddTask{
		Number1: rand.Intn(999),
		Number2: rand.Intn(999),
	}
	body, err := json.Marshal(addTask)
	go_rabbit.HandleError(err, "Error encoding JSON")

	//Prepare message struct to be send to queue
	//It has to be an instance of the aqmp publishing struct
	message := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	}

	//Send message through channel
	err = amqpChannel.Publish("calc", "", false, false, message)
	go_rabbit.HandleError(err, "Error publishing message")

	log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)
}
