package events

type EventListener interface {
	Listen(eventChannel chan map[string]string) error
}
