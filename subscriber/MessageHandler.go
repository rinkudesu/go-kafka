package subscriber

type MessageHandler interface {
	HandleMessage(message []byte)
	GetTopic() string
}
