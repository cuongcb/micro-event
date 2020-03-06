package amqp

import (
	"encoding/json"

	"github.com/cuongcb/micro-event/svcs/msgqueue"
	"github.com/streadway/amqp"
)

const (
	exchangeName = "events"
	exchangeType = "topic"
	durable      = true
	autoDelete   = false
	internal     = false
	noWait       = false
)

type eventEmitter struct {
	connection *amqp.Connection
}

// NewEventEmitter ...
func NewEventEmitter(conn *amqp.Connection) (msgqueue.Emitter, error) {
	emitter := &eventEmitter{
		connection: conn,
	}

	err := emitter.init()
	if err != nil {
		return nil, err
	}

	return emitter, nil
}

func (e *eventEmitter) init() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return channel.ExchangeDeclare(exchangeName, exchangeType, durable, autoDelete, internal, noWait, nil)
}

// Emit ...
func (e *eventEmitter) Emit(event msgqueue.Event) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event.EventName()},
		Body:        jsonEvent,
		ContentType: "application/json",
	}

	return channel.Publish(
		"events",          // exchangeName
		event.EventName(), // message routing key
		false,             // mandatory
		false,             // immediate
		msg)               // message
}
