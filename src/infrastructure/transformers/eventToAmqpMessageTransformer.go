package transformers

import (
	"encoding/json"
	"go-iam/src/domain/events"

	"github.com/streadway/amqp"
)

type EventToAMQPMessageTransformer struct{}

func (transformer *EventToAMQPMessageTransformer) Transform(event events.Event) (*amqp.Publishing, error) {
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        jsonBytes,
	}
	return &message, nil
}

func NewEventToAMQPMessageTransformer() *EventToAMQPMessageTransformer {
	transformer := EventToAMQPMessageTransformer{}
	return &transformer
}
