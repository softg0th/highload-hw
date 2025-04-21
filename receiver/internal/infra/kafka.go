package infra

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewProducer(bootstrapServers string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		return nil, err
	}

	admin, _ := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	admin.CreateTopics(context.Background(), []kafka.TopicSpecification{
		{
			Topic:             "test",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	})
	admin.Close()

	return &KafkaProducer{producer: p}, nil
}

func (kp *KafkaProducer) SendMessage(topic string, value []byte) error {
	log.Printf("sent")
	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}

func (kp *KafkaProducer) Close() error {
	if kp.producer != nil {
		kp.producer.Close()
	}
	return nil
}
