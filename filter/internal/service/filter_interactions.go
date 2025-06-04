package service

import (
	"context"
	"encoding/json"
	"filter/internal/core"
	"filter/internal/entities"
	"filter/internal/infra"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func FilterIt(ctx context.Context, rawMessage []byte, metadata interface{}) error {
	var msg entities.Message
	err := json.Unmarshal(rawMessage, &msg)
	if err != nil {
		return err
	}

	svc, _ := metadata.(*infra.Infra)
	cached, err := svc.Redis.Get(msg.Text)

	if err == nil && cached == "spam" {
		log.Println("%s is cached spam", msg.Text)
		return nil
	}
	if err != nil && err != redis.Nil {
		log.Println("Redis error: %v", err)
	}

	isItSpam := core.RunFilterPipeline(msg.Text)

	if isItSpam {
		_ = svc.Redis.Set(msg.Text, "spam", time.Hour)

		jsonMsg, _ := json.Marshal(msg)
		infra.SpamMessagesTotal.Inc()
		core.AddSpamMessage(jsonMsg)
	} else {
		err = svc.Rpc.StreamRequest(msg)
		return err
	}
	return nil
}
