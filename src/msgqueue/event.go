package msgqueue

// Event ...
type Event interface {
	PartitionKey() string
	EventName() string
}
