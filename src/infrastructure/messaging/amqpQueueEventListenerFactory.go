package messaging

import (
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

func (factory *AMQPQueueEventListenerFactory) CreateListener(event events.Event, eventConsumer events.EventConsumer) (events.EventListener, error) {
	eventConsumerName := eventConsumer.ConsumerName()
	exchange, err := factory.amqpExchangeManager.GetExchangeForEvent(event)
	if err != nil {
		return nil, err
	}

	_, err = factory.amqpChannel.QueueDeclare(
		eventConsumerName,
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
		eventConsumerName,
		eventConsumerName,
		*exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &AMQPQueueEventListener{
		amqpChannel:                  factory.amqpChannel,
		eventQueueName:               eventConsumerName,
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
