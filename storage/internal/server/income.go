package server

import (
	"io"
	"storage/internal/entities"
	storage "storage/internal/proto"
)

func (s *Server) StreamWithAck(stream storage.StorageService_StreamWithAckServer) error {
	for {
		in, err := stream.Recv()

		s.Logger.Info(map[string]interface{}{
			"message": "Got something interesting",
			"error":   false,
		})

		if err != nil {
			s.Logger.Info(map[string]interface{}{
				"message": "Error receiving message",
				"error":   true,
				"err":     err.Error(),
			})

			if err == io.EOF {
				s.Logger.Error(map[string]interface{}{
					"message": "End of stream",
					"error":   true,
				})
				return nil
			}
			return err
		}
		if in == nil {
			s.Logger.Info(map[string]interface{}{
				"message": "Received nil input",
				"error":   false,
			})
			continue
		}

		newEntity := entities.NewDocument(in.UserId, in.Timestamp.AsTime(), in.Text)

		err = s.Infra.Pool.Enqueue(newEntity)
		if err != nil {
			s.Logger.Error(map[string]interface{}{
				"message": "End of stream",
				"error":   true,
				"err":     err.Error(),
			})
		}
	}
}
