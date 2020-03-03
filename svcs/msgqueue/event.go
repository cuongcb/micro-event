package msgqueue

// Event ...
type Event interface {
	EventName() string
}
