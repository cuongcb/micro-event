package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/cuongcb/micro-event/msgqueue"
	"net/mail"
)

type eventEmitter struct {
	producer sarama.SyncProducer
}

func NewEventEmitter(client sarama.Client) (msgqueue.Listener, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := &eventEmitter{producer: producer}

	return emitter, nil
}

func (e *eventEmitter) Emit(event msgqueue.Event) error {
	envelope := messageEnvelope{
		EventName: event.EventName(),
		Payload:   event,
	}
	jsonBody, err := json.Marshal(&envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: event.EventName(),
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = e.producer.SendMessage(msg)
	return err
}
