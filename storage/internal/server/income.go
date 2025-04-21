package server

import (
	"io"
	"log"
	"storage/internal/entities"
	storage "storage/internal/proto"
)

func (s *Server) StreamWithAck(stream storage.StorageService_StreamWithAckServer) error {
	for {
		in, err := stream.Recv()
		log.Printf("got something interesting")
		if err != nil {
			log.Printf("error receiving message: %v", err)
			if err == io.EOF {
				log.Println("End of stream")
				return nil
			}
			return err
		}
		if in == nil {
			log.Println("Received nil input")
			continue
		}

		newEntity := entities.NewDocument(in.UserId, in.Timestamp.AsTime(), in.Text)

		err = s.Infra.Pool.Enqueue(newEntity)
		if err != nil {
			log.Printf("failed to enqueue: %v", err)
		}
	}
}
