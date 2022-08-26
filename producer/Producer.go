package producer

type Producer interface {
	Produce(topic string, message any) error
	Close()
}
