package service

func (s *Service) RunLoop() {
	for {
		msg, err := s.infra.Kc.ReadMessage(100)
		if err != nil {
			s.Logger.Error(map[string]interface{}{
				"message":           "Error reading message",
				"error":             true,
				"error_description": err.Error(),
			})
			continue
		}

		s.Logger.Info(map[string]interface{}{
			"message": "Received message",
			"error":   false,
			"msg":     msg.Value,
		})

		err = s.processingPool.Enqueue(msg.Value)
		if err != nil {
			s.Logger.Error(map[string]interface{}{
				"message":           "Error sending message",
				"error":             true,
				"error_description": err.Error(),
			})
			continue
		}
	}
}
