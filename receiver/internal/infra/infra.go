package infra

type Infra struct {
	Producer *KafkaProducer
}

func NewInfra(kafkaConn string) (*Infra, error) {
	producer, err := NewProducer(kafkaConn)
	if err != nil {
		return nil, err
	}
	return &Infra{Producer: producer}, nil
}
