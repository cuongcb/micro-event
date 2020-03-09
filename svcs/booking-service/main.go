package main

import (
	"log"

	queueAmqp "github.com/cuongcb/micro-event/svcs/msgqueue/amqp"
	"github.com/streadway/amqp"
)

func main() {
	log.SetPrefix("[booking-service] ")
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	listener, err := queueAmqp.NewEventListener(connection, "booking_queue")
	if err != nil {
		panic(err)
	}

	events, errors, err := listener.Listen("event.created")
	if err != nil {
		panic(err)
	}

	log.Println("Start listening...")

	for {
		select {
		case e := <-events:
			log.Println("receiving event:", e)
		case err := <-errors:
			log.Println("error:", err)
		}
	}

	// for {
	// 	select {
	// 	case e := <-events:
	// 		log.Println(e)
	// 	case err := <-errors:
	// 		log.Println(err)
	// 	}
	// }
}
