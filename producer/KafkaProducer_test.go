package producer

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/rinkudesu/go-kafka/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testMockMessage struct {
	Value string
}

func TestKafkaProducer_Produce(t *testing.T) {
	producer, testConsumer, ids := testGetProducer()
	defer producer.Close()
	_ = testConsumer.Subscribe(ids.String(), nil)
	defer func(testConsumer *kafka.Consumer) {
		_ = testConsumer.Close()
	}(testConsumer)
	msgChan := make(chan string)
	go func(msgStringChan chan string, test *testing.T) {
		readMessage, readErr := testConsumer.ReadMessage(time.Minute)
		assert.Nil(t, readErr)
		var readMessageObject testMockMessage
		_ = json.Unmarshal(readMessage.Value, &readMessageObject)
		msgStringChan <- readMessageObject.Value
		close(msgStringChan)
	}(msgChan, t)
	time.Sleep(time.Second)
	message := testMockMessage{Value: "test"}

	err := producer.Produce(ids.String(), message)

	assert.Nil(t, err)
	msgValue := <-msgChan
	assert.Equal(t, "test", msgValue)
}

func testGetProducer() (*KafkaProducer, *kafka.Consumer, uuid.UUID) {
	ids, _ := uuid.NewUUID()
	config := configuration.NewKafkaConfiguration("127.0.0.1:9092", "", "", ids.String(), ids.String())
	producer, _ := NewKafkaProducer(config)

	consumerConfig, _ := config.GetConsumerConfig()
	consumer, _ := kafka.NewConsumer(consumerConfig)
	admin, _ := kafka.NewAdminClientFromConsumer(consumer)
	_, _ = admin.CreateTopics(context.Background(), []kafka.TopicSpecification{{
		Topic:         ids.String(),
		NumPartitions: 1,
	}})
	return producer, consumer, ids
}
