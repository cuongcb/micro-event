package msgqueue

// Listener ...
type Listener interface {
	Listen(eventNames ...string) (<-chan Event, <-chan error, error)
}
