package producer

// Producer is responsible for sending messages over to the message broker
type Producer interface {
	// Produce sends a message to the broker with the given topic
	Produce(topic string, message any) error
	// Close closes connection and frees any used resources
	Close()
}
