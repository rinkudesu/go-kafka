package subscriber

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/rinkudesu/go-kafka/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockMessageHandler struct {
	topic      string
	Handled    int
	failHandle bool
}

func TestKafkaSubscriber_Subscribe(t *testing.T) {
	subscriber, producer, id := testGetSubscriber()
	defer producer.Close()
	mockHandler := &mockMessageHandler{topic: id.String()}
	err := subscriber.Subscribe(mockHandler)
	assert.Nil(t, err)

	err = subscriber.BeginHandle()
	time.Sleep(time.Second * 10)
	assert.Nil(t, err)
	topic := id.String()
	err = producer.Produce(&kafka.Message{Value: []byte{1, 2, 3}, TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}}, nil)
	assert.Nil(t, err)
	time.Sleep(time.Second * 10)

	assert.Equal(t, 1, mockHandler.Handled)
	subscriber.StopHandle()
	err = subscriber.Unsubscribe()
	assert.Nil(t, err)
	err = subscriber.Close()
	assert.Nil(t, err)
}

func (m *mockMessageHandler) HandleMessage(message []byte) bool {
	if len(message) != 3 || message[0] != 1 || message[1] != 2 || message[2] != 3 {
		return false
	}
	m.Handled++
	return !m.failHandle
}

func (m *mockMessageHandler) GetTopic() string {
	return m.topic
}

func testGetSubscriber() (*KafkaSubscriber, *kafka.Producer, uuid.UUID) {
	ids, _ := uuid.NewUUID()
	config := configuration.NewKafkaConfiguration("127.0.0.1:9092", "", "", ids.String(), ids.String())
	subscriber, _ := NewKafkaSubscriber(config)

	producerConfig, _ := config.GetProducerConfig()
	producer, _ := kafka.NewProducer(producerConfig)
	return subscriber, producer, ids
}
