package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/cuongcb/micro-event/msgqueue"
)

type eventListener struct {
	consumer sarama.Consumer
	partitions []int32
}

func NewEventListenr(client sarama.Client, partitions []int32) (msgqueue.Listener, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	listener := &eventListener{
		consumer:   consumer,
		partitions: partitions,
	}

	return listener, nil
}

func (e *eventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error) {

}