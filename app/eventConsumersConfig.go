package app

import (
	"fmt"
	"go-iam/src/application/sendPasswordResetToken"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func RunEventConsumers(container *dig.Container) {
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(consumer *sendPasswordResetToken.UserPasswordResetRequestedConsumer) {
			if err := consumer.Consume(); err != nil {
				logger.Fatal(fmt.Sprintf("Error running UserPasswordResetRequestedConsumer: %s", err.Error()))
			}
		}), logger)
	}); err != nil {
		panic(fmt.Sprintf("Error adding event consumers to the dependency container %s", err.Error()))
	}
}
