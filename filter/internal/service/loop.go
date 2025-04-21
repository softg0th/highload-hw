package service

import "log"

func (s *Service) RunLoop() {
	for {
		msg, err := s.infra.Kc.ReadMessage(100)
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}
		log.Printf("received message: %s\n", msg.Value)
		err = s.processingPool.Enqueue(msg.Value)
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}
	}
}
