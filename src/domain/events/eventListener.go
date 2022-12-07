package events

type EventListener interface {
	Listen(eventChannel chan map[string]interface{}) error
}
