package transformers

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type EventToAMQPMessageTransformer struct{}

func (transformer *EventToAMQPMessageTransformer) Transform(dto interface{}) (*amqp.Publishing, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonBytes,
	}
	return &message, nil
}

func NewEventToAMQPMessageTransformer() *EventToAMQPMessageTransformer {
	transformer := EventToAMQPMessageTransformer{}
	return &transformer
}
