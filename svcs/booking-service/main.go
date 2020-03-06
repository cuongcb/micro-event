package main

import (
	queueAmqp "github.com/cuongcb/micro-event/svcs/msgqueue/amqp"
	"github.com/streadway/amqp"
)

func main() {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	emitter, err := queueAmqp.NewEventEmitter(connection)
	if err != nil {
		panic(err)
	}
}
