package subscriber

// MessageHandler is responsible for processing messages read from the message broker.
type MessageHandler interface {
	// HandleMessage parses the message received from the broker and takes any necessary action to process it.
	// True is returned if the message has been correctly handled, false otherwise.
	HandleMessage(message []byte) bool
	// GetTopic returns a topic, messages from which can be handled by this handler
	GetTopic() string
}
