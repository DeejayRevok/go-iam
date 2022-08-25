package sendPasswordResetToken

import (
	"fmt"
	"go-uaa/src/domain/events"
	"go-uaa/src/domain/user"
	"reflect"

	"go.uber.org/zap"
)

type UserPasswordResetRequestedConsumer struct {
	eventListener         events.EventListener
	sendResetTokenUseCase *SendPasswordResetTokenUseCase
	logger                *zap.Logger
}

func (consumer *UserPasswordResetRequestedConsumer) Consume() error {
	eventMapChannel := make(chan map[string]string)
	err := consumer.eventListener.Listen(eventMapChannel)
	if err != nil {
		return err
	}
	go consumer.consumeEventMaps(eventMapChannel)
	return nil
}

func (consumer *UserPasswordResetRequestedConsumer) consumeEventMaps(eventMapChannel chan map[string]string) {
	for eventMap := range eventMapChannel {
		err := consumer.handleEventMap(eventMap)
		if err != nil {
			consumer.logger.Warn(fmt.Sprintf("Error handling event: %s", err.Error()))
		}
	}
}

func (consumer *UserPasswordResetRequestedConsumer) handleEventMap(eventMap map[string]string) error {
	event := user.UserPasswordResetRequestedEventFromMap(eventMap)
	sendResetTokenRequest := SendPasswordResetTokenRequest{
		UserID:     event.UserID,
		ResetToken: event.ResetToken,
	}
	if useCaseResponse := consumer.sendResetTokenUseCase.Execute(&sendResetTokenRequest); useCaseResponse.Err != nil {
		return useCaseResponse.Err
	}
	return nil
}

func NewUserPasswordResetRequestedConsumer(eventListenerFactory events.EventListenerFactory, useCase *SendPasswordResetTokenUseCase, logger *zap.Logger) *UserPasswordResetRequestedConsumer {
	eventName := reflect.TypeOf(user.UserPasswordResetRequestedEvent{}).Name()
	eventListener, err := eventListenerFactory.CreateListener(eventName)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating listener %s: %s", eventName, err.Error()))
	}
	return &UserPasswordResetRequestedConsumer{
		eventListener:         eventListener,
		sendResetTokenUseCase: useCase,
		logger:                logger,
	}
}
