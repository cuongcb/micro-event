package amqp

import (
	"encoding/json"

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

// EventEmitter ...
type EventEmitter struct {
	connection *amqp.Connection
}

// NewEventEmitter ...
func NewEventEmitter(conn *amqp.Connection) (*EventEmitter, error) {
	emitter := &EventEmitter{
		connection: conn,
	}

	err := emitter.init()
	if err != nil {
		return nil, err
	}

	return emitter, nil
}

func (e *EventEmitter) init() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	return channel.ExchangeDeclare(exchangeName, exchangeType, durable, autoDelete, internal, noWait, nil)
}

func (e *EventEmitter) Emit(event Event) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	chan, err := e.connection.Channel();
	if err != nil {
		return err
	}

	defer chan.Close()

	msg := amqp.Publishing{
		Headers: amqpTable{"x-event-name": event.EventName()},
		Body: jsonEvent,
		ContentType: "application/json",
	}

	return chan.Publish(
		"events", // exchangeName
		event.EventName(), // message routing key
		false, // mandatory
		false, // immediate
		msg // message
	)
}
