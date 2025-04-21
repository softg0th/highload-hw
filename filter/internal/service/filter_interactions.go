package service

import (
	"context"
	"encoding/json"
	"filter/internal/core"
	"filter/internal/entities"
	"filter/internal/infra"
	"log"
)

func FilterIt(ctx context.Context, rawMessage []byte, metadata interface{}) error {
	var msg entities.Message

	err := json.Unmarshal(rawMessage, &msg)

	if err != nil {
		return err
	}

	isItSpam := core.RunFilterPipeline(msg.Text)

	if isItSpam {
		jsonMsg, _ := json.Marshal(msg)
		log.Printf("spam: %s\n", msg.Text)
		core.AddSpamMessage(jsonMsg)
	} else {
		log.Printf("not spam: %s\n", msg.Text)
		svc, _ := metadata.(*infra.Infra)
		err = svc.Rpc.StreamRequest(msg)
		log.Printf("error: %s\n", err)
		return err
	}
	return nil
}
