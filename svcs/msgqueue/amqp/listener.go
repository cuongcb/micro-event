package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type EventListener struct {
	connection *amqp.Connection
	queue      string
}

func NewEventListener(conn *amqp.Connection, queue string) (*EventListener, error) {
	listener := &EventListener{
		connection: conn,
		queue:      queue,
	}

	err := listener.init()
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (e *EventListener) init() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	_, err := channel.QueueDeclare(e.queue, true, false, false, false, nil)

	return err
}

func (e *EventListener) Listen(eventNames ...string) (<-chan Event, <-chan error, error) {
	channel, err := e.connection.Channel()
	if err != nil {
		return nil, nil err
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

	events := make(chan Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("no x-event-name header found in msg")
				msg.Nack(false)
				continue
			}

			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name header is not string, but %t", rawEventName)
				msg.Nack(false)
				continue
			}

			var event Event
			switch eventName {
			case "event.created":
				event = new(contracts.EventCreatedEvent)
			case default:
				errors <- fmt.Errorf("event type %s is unknown", eventName)
				continue // iterate new msg
			}

			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- err
				continue
			}

			events <- event
		}
	}()

	return events, errors, nil
}
