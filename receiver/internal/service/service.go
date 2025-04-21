package service

type Service struct {
	Sender *Sender
}

func NewService(sender *Sender) *Service {
	return &Service{
		Sender: sender,
	}
}
