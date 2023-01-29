package messaging

import (
	"go-iam/src/domain/events"
	"go-iam/src/infrastructure/transformers"

	"github.com/streadway/amqp"
)

type AMQPExchangeEventPublisher struct {
	amqpChannel             *amqp.Channel
	amqpExchangeManager     *AMQPExchangeManager
	eventMessageTransformer *transformers.EventToAMQPMessageTransformer
}

func (publisher *AMQPExchangeEventPublisher) Publish(event events.Event) error {
	exchange, err := publisher.amqpExchangeManager.GetExchangeForEvent(event)
	if err != nil {
		return err
	}
	message, err := publisher.eventMessageTransformer.Transform(event)
	if err != nil {
		return err
	}
	if err := publisher.amqpChannel.Publish(*exchange, "", false, false, *message); err != nil {
		return err
	}

	return nil
}

func NewAMQPExchangeEventPublisher(amqpChannel *amqp.Channel, amqpExchangeManager *AMQPExchangeManager, eventTransformer *transformers.EventToAMQPMessageTransformer) *AMQPExchangeEventPublisher {
	publisher := AMQPExchangeEventPublisher{
		amqpChannel:             amqpChannel,
		amqpExchangeManager:     amqpExchangeManager,
		eventMessageTransformer: eventTransformer,
	}
	return &publisher
}
