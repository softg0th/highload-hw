package service

import (
	"context"
	"encoding/json"
	"filter/internal/core"
	"filter/internal/entities"
	"filter/internal/infra"
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
		infra.SpamMessagesTotal.Inc()
		core.AddSpamMessage(jsonMsg)
	} else {
		svc, _ := metadata.(*infra.Infra)
		err = svc.Rpc.StreamRequest(msg)
		return err
	}
	return nil
}
