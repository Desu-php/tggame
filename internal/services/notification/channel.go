package notification

type Channel interface {
	Send(message string, receiver *Receiver) error
}

type Receiver struct {
	ID string
}
