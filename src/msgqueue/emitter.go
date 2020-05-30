package msgqueue

// Emitter ...
type Emitter interface {
	Emit(e Event) error
}
