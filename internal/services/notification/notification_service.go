package notification

import "fmt"

type Service struct {
	channels []Channel
}

func NewService(channels []Channel) *Service {
	return &Service{
		channels,
	}
}

func (s *Service) NotifyAll(message string, receiver *Receiver) error {
	var errs []error

	for _, channel := range s.channels {
		if err := channel.Send(message, receiver); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("one or more channels failed: %v", errs)
	}

	return nil
}
