package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/cuongcb/micro-event/msgqueue"
	"github.com/cuongcb/micro-event/msgqueue/contracts"
	"github.com/streadway/amqp"
)

type eventListener struct {
	connection *amqp.Connection
	queue      string
}

// NewEventListener ...
func NewEventListener(conn *amqp.Connection, queue string) (msgqueue.Listener, error) {
	listener := &eventListener{
		connection: conn,
		queue:      queue,
	}

	err := listener.init()
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (e *eventListener) init() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	_, err = channel.QueueDeclare(e.queue, true, false, false, false, nil)

	return err
}

func (e *eventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := e.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	defer channel.Close()

	for _, eventName := range eventNames {
		if err := channel.QueueBind(e.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}

	msgs, err := channel.Consume(e.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("no x-event-name header found in msg")
				msg.Nack(false, false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name header is not string, but %t", rawEventName)
				msg.Nack(false, false)
				continue
			}

			var event msgqueue.Event
			switch eventName {
			case "event.created":
				event = new(contracts.EventCreatedEvent)
			default:
				msg.Nack(false, false)
				errors <- fmt.Errorf("event type %s is unknown", eventName)
				continue // iterate new msg
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				msg.Nack(false, false)
				errors <- err
				continue
			}

			msg.Ack(false)
			events <- event
		}
	}()

	return events, errors, nil
}
