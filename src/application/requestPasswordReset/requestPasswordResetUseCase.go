package requestPasswordReset

import (
	"context"
	"fmt"
	"go-iam/src/domain/events"
	"go-iam/src/domain/hash"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"time"

	"github.com/google/uuid"
)

type RequestPasswordResetUseCase struct {
	userRepository              user.UserRepository
	userPasswordResetRepository user.UserPasswordResetRepository
	hasher                      hash.Hasher
	eventPublisher              events.EventPublisher
	logger                      internals.Logger
}

func (useCase *RequestPasswordResetUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*RequestPasswordResetRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting password reset request for user with mail %s", validatedRequest.Email))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished password reset request for user with mail %s", validatedRequest.Email))

	requestUser, err := useCase.userRepository.FindByEmail(ctx, validatedRequest.Email)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	resetToken, resetTokenHash, err := useCase.getResetToken()
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	if err = useCase.deleteExistingResetToken(ctx, requestUser.ID); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	userPasswordReset := user.UserPasswordReset{
		Token:      resetTokenHash,
		Expiration: time.Now().Add(time.Minute * 15),
		UserID:     requestUser.ID,
	}
	if err = useCase.userPasswordResetRepository.Save(ctx, userPasswordReset); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if err = useCase.publishEvent(resetToken, requestUser); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (useCase *RequestPasswordResetUseCase) deleteExistingResetToken(ctx context.Context, userID uuid.UUID) error {
	userPasswordReset, err := useCase.userPasswordResetRepository.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if userPasswordReset != nil {
		err = useCase.userPasswordResetRepository.Delete(ctx, *userPasswordReset)
	}
	return err
}

func (useCase *RequestPasswordResetUseCase) getResetToken() (string, string, error) {
	resetToken := uuid.New().String()
	resetTokenHash, err := useCase.hasher.Hash(resetToken)
	if err != nil {
		return "", "", err
	}
	return resetToken, *resetTokenHash, nil
}

func (useCase *RequestPasswordResetUseCase) publishEvent(resetToken string, requester *user.User) error {
	event := user.UserPasswordResetRequestedEvent{
		ResetToken: resetToken,
		UserID:     requester.ID.String(),
	}
	return useCase.eventPublisher.Publish(&event)
}

func NewRequestPasswordResetUseCase(userRepository user.UserRepository, userPasswordResetRepository user.UserPasswordResetRepository, eventPublisher events.EventPublisher, hasher hash.Hasher, logger internals.Logger) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepository:              userRepository,
		userPasswordResetRepository: userPasswordResetRepository,
		hasher:                      hasher,
		eventPublisher:              eventPublisher,
		logger:                      logger,
	}
}
