package entities

type KafkaResult struct {
	ID      uint64
	Success bool
	Error   error
}

type KafkaTask struct {
	ID    uint64
	Msg   []byte
	Reply chan KafkaResult
}

func NewKafkaTask(id uint64, msg []byte) KafkaTask {
	return KafkaTask{
		ID:    id,
		Msg:   msg,
		Reply: make(chan KafkaResult, 1),
	}
}
