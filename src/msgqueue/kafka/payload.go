package kafka

type messageEnvelope struct {
	EventName string `json:"event_name"`
	Payload interface{} `json:"payload"`
}
