package messaging

import (
	"go-iam/src/domain/events"

	"github.com/streadway/amqp"
)

type AMQPExchangeManager struct {
	amqpChannel *amqp.Channel
	exchanges   map[string]string
}

func (manager *AMQPExchangeManager) GetExchangeForEvent(event events.Event) (*string, error) {
	eventType := event.EventName()
	exchange := manager.exchanges[eventType]
	if exchange != "" {
		return &exchange, nil
	}

	err := manager.createExchange(eventType)
	if err != nil {
		return nil, err
	}
	manager.exchanges[eventType] = eventType
	return &eventType, nil
}

func (manager *AMQPExchangeManager) createExchange(name string) error {
	return manager.amqpChannel.ExchangeDeclare(name, "fanout", true, false, false, false, nil)
}

func NewAMQPExchangeManager(amqpChannel *amqp.Channel) *AMQPExchangeManager {
	manager := AMQPExchangeManager{
		amqpChannel: amqpChannel,
		exchanges:   make(map[string]string),
	}
	return &manager
}
