package events

type EventConsumer interface {
	Consume() error
}
