package producer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"go-kafka/configuration"
	"time"
)

const (
	sendAttempts = 10
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(configuration *configuration.KafkaConfiguration) (*KafkaProducer, error) {
	producerConfig, err := configuration.GetProducerConfig()
	if err != nil {
		return nil, err
	}

	kafkaProducer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return nil, err
	}

	go logErrors(kafkaProducer)
	return &KafkaProducer{producer: kafkaProducer}, nil
}

func (producer *KafkaProducer) Produce(topic string, message any) error {
	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return producer.produceWithWait(topic, messageJson)
}

func (producer *KafkaProducer) produceWithWait(topic string, message []byte) error {
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}

	var err error
	for i := 0; i < sendAttempts; i++ {
		err = producer.producer.Produce(kafkaMessage, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				log.Warningf("Producer queue is full, waiting for a second and trying again, attempt %d", i)
				time.Sleep(time.Second)
				continue
			}
			log.Warningf("Message sending failed: %v", err)
		} else {
			break
		}
	}

	return err
}

func (producer *KafkaProducer) Close() {
	for producer.producer.Flush(1000) > 0 {
		log.Info("Still waiting to flush outstanding messages")
	}
	producer.producer.Close()
}

func logErrors(producer *kafka.Producer) {
	for e := range producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				log.Errorf("Failed to send Kafka message: %v", m.TopicPartition.Error)
			} else {
				log.Debugf("Delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		case kafka.Error:
			log.Warningf("Kafka error %s", ev.String())
		default:
			log.Infof("Ignored Kafka event: %s", ev)
		}
	}
}
