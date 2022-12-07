package messaging

import (
	"fmt"
	"go-iam/src/infrastructure/transformers"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type AMQPQueueEventListener struct {
	amqpChannel                  *amqp.Channel
	eventQueueName               string
	amqpDeliveryToMapTransformer *transformers.AMQPDeliveryToMapTransformer
	logger                       *zap.Logger
}

func (listener *AMQPQueueEventListener) Listen(eventChannel chan map[string]interface{}) error {
	consumerTag := fmt.Sprintf("%sConsumer.%s", listener.eventQueueName, uuid.New().String())
	deliveries, err := listener.amqpChannel.Consume(
		listener.eventQueueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}
	go listener.handleDelivery(deliveries, eventChannel)
	return nil
}

func (listener *AMQPQueueEventListener) handleDelivery(deliveries <-chan amqp.Delivery, eventChannel chan map[string]interface{}) {
	for delivery := range deliveries {
		listener.logger.Info(fmt.Sprintf("got %dB delivery: [%v] %q", len(delivery.Body), delivery.DeliveryTag, delivery.Body))
		eventMap, err := listener.amqpDeliveryToMapTransformer.Transform(&delivery)
		if err != nil {
			listener.logger.Warn(fmt.Sprintf("Error transforming message %s to map: %s", string(delivery.Body), err.Error()))
		}
		eventChannel <- eventMap
		if err = delivery.Ack(false); err != nil {
			listener.logger.Warn(fmt.Sprintf("Error acknowledging message %s due to %s", string(delivery.Body), err.Error()))
		}
	}
}
