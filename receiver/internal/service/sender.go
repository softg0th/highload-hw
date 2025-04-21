package service

import (
	"context"
	"encoding/json"
	"log"
	"receiver/internal/domain/entities"
	"receiver/internal/infra"
	"sync/atomic"
)

type Sender struct {
	pool *infra.Pool
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
		return entities.KafkaResult{}, err
	}
	log.Printf("executed")
	id := nextID()
	task := entities.NewKafkaTask(id, messageBytes)
	err = s.pool.Enqueue(task)
	if err != nil {
		log.Printf("failed to enqueue task: %v", err)
	}
	select {
	case <-ctx.Done():
		return entities.KafkaResult{}, ctx.Err()
	case res := <-task.Reply:
		return res, nil
	}
}

func NewSender(pool *infra.Pool) *Sender {
	return &Sender{
		pool: pool,
	}
}
