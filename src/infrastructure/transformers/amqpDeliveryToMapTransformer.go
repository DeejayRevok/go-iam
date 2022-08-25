package transformers

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type AMQPDeliveryToMapTransformer struct{}

func (*AMQPDeliveryToMapTransformer) Transform(delivery *amqp.Delivery) (map[string]string, error) {
	transformedDelivery := map[string]string{}
	err := json.Unmarshal(delivery.Body, &transformedDelivery)
	if err != nil {
		return nil, err
	}
	return transformedDelivery, nil
}

func NewAMQPDeliveryToMapTransformer() *AMQPDeliveryToMapTransformer {
	return &AMQPDeliveryToMapTransformer{}
}
