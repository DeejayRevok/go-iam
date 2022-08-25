package messaging

import (
	"errors"
	"time"

	"github.com/streadway/amqp"
)

type AMQPHealthChecker struct {
	amqpConnection  *amqp.Connection
	amqpHealthError error
}

func (checker *AMQPHealthChecker) Check() error {
	return checker.amqpHealthError
}

func (checker *AMQPHealthChecker) monitorHealth() {
	notify := checker.amqpConnection.NotifyClose(make(chan *amqp.Error))
	for {
		err := <-notify
		if err != nil {
			checker.amqpHealthError = errors.New(err.Error())
		} else {
			checker.amqpHealthError = errors.New("AMQP connection lost")
		}
		time.Sleep(1 * time.Second)
	}
}

func NewAMQPHealthChecker(amqpConnection *amqp.Connection) *AMQPHealthChecker {
	checker := AMQPHealthChecker{
		amqpConnection:  amqpConnection,
		amqpHealthError: nil,
	}
	go checker.monitorHealth()
	return &checker
}
