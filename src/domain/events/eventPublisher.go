package events

type EventPublisher interface {
	Publish(event interface{}) error
}
