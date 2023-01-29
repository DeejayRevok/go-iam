package events

type EventListenerFactory interface {
	CreateListener(event Event, eventConsumer EventConsumer) (EventListener, error)
}
