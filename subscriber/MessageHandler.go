package subscriber

type MessageHandler interface {
	HandleMessage(message []byte) bool
	GetTopic() string
}
