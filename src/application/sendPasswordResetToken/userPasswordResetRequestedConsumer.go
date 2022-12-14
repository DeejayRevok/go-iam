package sendPasswordResetToken

import (
	"context"
	"fmt"
	"go-iam/src/domain/events"
	"go-iam/src/domain/user"
	"reflect"

	"go.uber.org/zap"
)

type UserPasswordResetRequestedConsumer struct {
	eventListener         events.EventListener
	sendResetTokenUseCase *SendPasswordResetTokenUseCase
	logger                *zap.Logger
}

func (consumer *UserPasswordResetRequestedConsumer) Consume() error {
	eventMapChannel := make(chan map[string]interface{})
	err := consumer.eventListener.Listen(eventMapChannel)
	if err != nil {
		return err
	}
	go consumer.consumeEventMaps(eventMapChannel)
	return nil
}

func (consumer *UserPasswordResetRequestedConsumer) consumeEventMaps(eventMapChannel chan map[string]interface{}) {
	for eventMap := range eventMapChannel {
		err := consumer.handleEventMap(eventMap)
		if err != nil {
			consumer.logger.Warn(fmt.Sprintf("Error handling event: %s", err.Error()))
		}
	}
}

func (consumer *UserPasswordResetRequestedConsumer) handleEventMap(eventMap map[string]interface{}) error {
	event := user.UserPasswordResetRequestedEventFromMap(eventMap)
	sendResetTokenRequest := SendPasswordResetTokenRequest{
		UserID:     event.UserID,
		ResetToken: event.ResetToken,
	}
	context := context.Background()
	if useCaseResponse := consumer.sendResetTokenUseCase.Execute(context, &sendResetTokenRequest); useCaseResponse.Err != nil {
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
