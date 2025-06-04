package infra

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewKafkaConsumer(bootstrapServers string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{"test"}, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c}, nil
}

func (kc *KafkaConsumer) ReadMessage(timeoutMs int) (*kafka.Message, error) {
	msg, err := kc.consumer.ReadMessage(time.Second * time.Duration(timeoutMs))
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (kc *KafkaConsumer) Close() error {
	if kc.consumer != nil {
		return kc.consumer.Close()
	}
	return nil
}
