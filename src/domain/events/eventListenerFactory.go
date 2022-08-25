package events

type EventListenerFactory interface {
	CreateListener(eventName string) (EventListener, error)
}
