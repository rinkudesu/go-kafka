package configuration

import (
	"errors"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	EnvValueMissingErr = errors.New("required configuration value is missing from env variables")
)

// KafkaConfiguration holds data allowing for broker connection
type KafkaConfiguration struct {
	serverAddress   string
	user            string
	password        string
	clientId        string
	consumerGroupId string
}

// NewKafkaConfiguration creates new KafkaConfiguration using supplied data
func NewKafkaConfiguration(serverAddress string, user string, password string, clientId string, consumerGroupId string) *KafkaConfiguration {
	return &KafkaConfiguration{serverAddress: serverAddress, user: user, password: password, clientId: clientId, consumerGroupId: consumerGroupId}
}

// NewKafkaConfigurationFromEnv loads all necessary data to create KafkaConfiguration from env variables
func NewKafkaConfigurationFromEnv() (*KafkaConfiguration, error) {
	serverAddress := os.Getenv("RINKU_KAFKA_ADDRESS")
	user := os.Getenv("RINKU_KAFKA_USER")
	password := os.Getenv("RINKU_KAFKA_PASSWORD")
	clientId := os.Getenv("RINKU_KAFKA_CLIENT_ID")
	consumerGroupId := os.Getenv("RINKU_KAFKA_CONSUMER_GROUP_ID")

	if serverAddress == "" || clientId == "" || consumerGroupId == "" {
		return nil, EnvValueMissingErr
	}
	return NewKafkaConfiguration(serverAddress, user, password, clientId, consumerGroupId), nil
}

func (configuration *KafkaConfiguration) fillSasl(configMap *kafka.ConfigMap) error {
	if configuration.user == "" || configuration.password == "" {
		return nil
	}

	if err := configMap.SetKey("sasl.mechanisms", "PLAIN"); err != nil {
		return err
	}
	if err := configMap.SetKey("security.protocol", "PLAIN"); err != nil {
		return err
	}
	if err := configMap.SetKey("sasl.username", configuration.user); err != nil {
		return err
	}
	if err := configMap.SetKey("sasl.password", configuration.password); err != nil {
		return err
	}
	return nil
}

func (configuration *KafkaConfiguration) getCommon() (*kafka.ConfigMap, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": configuration.serverAddress,
		"client.id":         configuration.clientId,
	}
	err := configuration.fillSasl(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (configuration *KafkaConfiguration) GetConsumerConfig() (*kafka.ConfigMap, error) {
	config, err := configuration.getCommon()
	if err != nil {
		return nil, err
	}
	if err := config.SetKey("enable.auto.commit", false); err != nil {
		return nil, err
	}
	if err := config.SetKey("auto.offset.reset", "earliest"); err != nil {
		return nil, err
	}
	if err := config.SetKey("group.id", configuration.consumerGroupId); err != nil {
		return nil, err
	}
	if err := config.SetKey("partition.assignment.strategy", "cooperative-sticky"); err != nil {
		return nil, err
	}
	return config, nil
}

func (configuration *KafkaConfiguration) GetProducerConfig() (*kafka.ConfigMap, error) {
	return configuration.getCommon()
}
