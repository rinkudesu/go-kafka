package subscriber

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rinkudesu/go-kafka/configuration"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	NotSubscribedErr   = errors.New("subscription is not yet active")
	AlreadyHandlingErr = errors.New("this subscriber is already handling messages")
)

type KafkaSubscriber struct {
	subscriber *kafka.Consumer
	handler    MessageHandler

	terminate        bool //todo: this is not the best, but it'll do for now - do some proper cancellation thing here
	handlerWaitGroup *sync.WaitGroup
}

func NewKafkaSubscriber(configuration *configuration.KafkaConfiguration) (*KafkaSubscriber, error) {
	subscriberConfig, err := configuration.GetConsumerConfig()
	if err != nil {
		return nil, err
	}

	subscriber, err := kafka.NewConsumer(subscriberConfig)
	if err != nil {
		return nil, err
	}
	return &KafkaSubscriber{subscriber: subscriber}, nil
}

func (subscriber *KafkaSubscriber) Subscribe(handler MessageHandler) error {
	subscriber.handler = handler

	//ensure topic exists manually due to bug: https://github.com/confluentinc/confluent-kafka-go/issues/788
	adminClient, err := kafka.NewAdminClientFromConsumer(subscriber.subscriber)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	// ignore errors here, as fails will appear later anyway
	_, _ = adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{{
		Topic:         handler.GetTopic(),
		NumPartitions: 1,
	}})

	err = subscriber.subscriber.Subscribe(handler.GetTopic(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (subscriber *KafkaSubscriber) Unsubscribe() error {
	if subscriber.handler == nil {
		return nil
	}

	err := subscriber.subscriber.Unsubscribe()
	if err != nil {
		return err
	}
	subscriber.handler = nil
	return nil
}

func (subscriber *KafkaSubscriber) BeginHandle() error {
	if subscriber.handler == nil {
		return NotSubscribedErr
	}
	if subscriber.handlerWaitGroup != nil {
		return AlreadyHandlingErr
	}

	subscriber.terminate = false
	subscriber.handlerWaitGroup = &sync.WaitGroup{}
	subscriber.handlerWaitGroup.Add(1)
	go listenForMessages(subscriber)
	return nil
}

func (subscriber *KafkaSubscriber) StopHandle() {
	if subscriber.handlerWaitGroup == nil {
		return
	}

	subscriber.terminate = true
	subscriber.handlerWaitGroup.Wait()
	subscriber.handlerWaitGroup = nil
}

func (subscriber *KafkaSubscriber) Close() error {
	subscriber.StopHandle()
	if err := subscriber.Unsubscribe(); err != nil {
		return err
	}
	if err := subscriber.subscriber.Close(); err != nil {
		return err
	}
	return nil
}

func listenForMessages(subscriber *KafkaSubscriber) {
	defer subscriber.handlerWaitGroup.Done()
	for !(subscriber.terminate) {
		message, err := subscriber.subscriber.ReadMessage(time.Second * 10)
		if err == nil {
			result := subscriber.handler.HandleMessage(message.Value)
			if result {
				_, err := subscriber.subscriber.CommitMessage(message)
				if err != nil {
					log.Warningf("Failed to commit message: %s", err)
				}
			}
		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			log.Debug("Timed out while waiting for kafka message")
		} else {
			log.Warningf("An error has occured while waiting for kafka message: %s", err.Error())
		}
	}
}
