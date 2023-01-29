package events

type EventPublisher interface {
	Publish(event Event) error
}
