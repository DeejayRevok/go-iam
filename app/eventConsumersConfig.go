package app

import (
	"fmt"
	"go-uaa/src/application/sendPasswordResetToken"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func RunEventConsumers(container *dig.Container) {
	container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(consumer *sendPasswordResetToken.UserPasswordResetRequestedConsumer) {
			if err := consumer.Consume(); err != nil {
				panic(fmt.Sprintf("Error running UserPasswordResetRequestedConsumer: %s", err.Error()))
			}
		}), logger)
	})
}
