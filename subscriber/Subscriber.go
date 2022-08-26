package subscriber

// Subscriber is a consumer for a message broker. It will listen for new messages from the broker and pass them for handling.
type Subscriber interface {
	// Subscribe subscribed to a topic defined in handler.
	Subscribe(handler MessageHandler) error
	// Unsubscribe removes current subscription.
	Unsubscribe() error

	// BeginHandle starts handling messages received from the broker.
	// This call starts a goroutine in the background and will not block.
	BeginHandle() error
	// StopHandle stops message handling and blocks until current message processing finishes.
	StopHandle()

	// Close disables handling, unsubscribed, and disconnects from the broker.
	Close() error
}
