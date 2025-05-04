package service

import (
	"context"
	"encoding/json"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"receiver/internal/domain/entities"
	"receiver/internal/infra"
	"sync/atomic"
)

type Sender struct {
	pool   *infra.Pool
	logger *logstash_logger.Logstash
}

var globalID uint64

func nextID() uint64 {
	return atomic.AddUint64(&globalID, 1)
}

func transformMessage(message entities.Message) ([]byte, error) {
	srMsg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return srMsg, nil
}

func (s *Sender) Send(ctx context.Context, message entities.Message) (entities.KafkaResult, error) {
	messageBytes, err := transformMessage(message)

	if err != nil {
		s.logger.Error(map[string]interface{}{
			"message": "failed to transform message",
			"error":   true,
		})
		return entities.KafkaResult{}, err
	}

	s.logger.Info(map[string]interface{}{
		"message": "Incoming message transformed",
		"error":   false,
	})

	id := nextID()
	task := entities.NewKafkaTask(id, messageBytes)
	err = s.pool.Enqueue(task)

	if err != nil {
		s.logger.Error(map[string]interface{}{
			"message": "failed to enqueue task",
			"error":   true,
		})
		return entities.KafkaResult{}, err
	}

	select {
	case <-ctx.Done():
		s.logger.Error(map[string]interface{}{
			"message": "context cancelled",
		})
		return entities.KafkaResult{}, ctx.Err()
	case res := <-task.Reply:
		s.logger.Info(map[string]interface{}{
			"message": "Received reply from task",
			"error":   false,
			"id":      id,
		})
		return res, nil
	}
}

func NewSender(pool *infra.Pool, logger *logstash_logger.Logstash) *Sender {
	return &Sender{
		pool:   pool,
		logger: logger,
	}
}
