package events

type Event interface{
	EventName() string
}