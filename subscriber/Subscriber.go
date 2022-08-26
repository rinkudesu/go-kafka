package subscriber

type Subscriber interface {
	Subscribe(handler MessageHandler) error
	Unsubscribe() error

	BeginHandle() error
	StopHandle()

	Close() error
}
