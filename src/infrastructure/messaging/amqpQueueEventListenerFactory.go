package messaging

import (
	"fmt"
	"go-iam/src/domain/events"
	"go-iam/src/infrastructure/transformers"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPQueueEventListenerFactory struct {
	amqpChannel                  *amqp.Channel
	amqpExchangeManager          *AMQPExchangeManager
	amqpDeliveryToMapTransformer *transformers.AMQPDeliveryToMapTransformer
	logger                       *zap.Logger
}

func (factory *AMQPQueueEventListenerFactory) CreateListener(eventName string) (events.EventListener, error) {
	eventQueueName := fmt.Sprintf("iam.%s", eventName)
	exchange, err := factory.amqpExchangeManager.GetExchangeForEvent(eventName)
	if err != nil {
		return nil, err
	}

	_, err = factory.amqpChannel.QueueDeclare(
		eventQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = factory.amqpChannel.QueueBind(
		eventQueueName,
		"iam",
		*exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &AMQPQueueEventListener{
		amqpChannel:                  factory.amqpChannel,
		eventQueueName:               eventQueueName,
		amqpDeliveryToMapTransformer: factory.amqpDeliveryToMapTransformer,
		logger:                       factory.logger,
	}, nil
}

func NewAMQPQueueEventListenerFactory(amqpChannel *amqp.Channel, amqpExchangeManager *AMQPExchangeManager, amqpDeliveryToMapTransformer *transformers.AMQPDeliveryToMapTransformer, logger *zap.Logger) *AMQPQueueEventListenerFactory {
	return &AMQPQueueEventListenerFactory{
		amqpChannel:                  amqpChannel,
		amqpExchangeManager:          amqpExchangeManager,
		amqpDeliveryToMapTransformer: amqpDeliveryToMapTransformer,
		logger:                       logger,
	}
}
