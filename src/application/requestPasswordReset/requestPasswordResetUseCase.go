package requestPasswordReset

import (
	"fmt"
	"go-uaa/src/domain/events"
	"go-uaa/src/domain/hash"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestPasswordResetUseCase struct {
	userRepository              user.UserRepository
	userPasswordResetRepository user.UserPasswordResetRepository
	hasher                      hash.Hasher
	eventPublisher              events.EventPublisher
	logger                      *zap.Logger
}

func (useCase *RequestPasswordResetUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*RequestPasswordResetRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting password reset request for user with mail %s", validatedRequest.Email))
	defer useCase.logger.Info(fmt.Sprintf("Finished password reset request for user with mail %s", validatedRequest.Email))

	requestUser, err := useCase.userRepository.FindByEmail(validatedRequest.Email)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	resetToken, resetTokenHash, err := useCase.getResetToken()
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}

	userPasswordReset := user.UserPasswordReset{
		Token:      resetTokenHash,
		Expiration: time.Now().Add(time.Minute * 15),
		UserID:     requestUser.ID,
	}
	if err = useCase.userPasswordResetRepository.Save(userPasswordReset); err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	useCase.publishEvent(resetToken, requestUser)
	return internals.EmptyUseCaseResponse()
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
	return useCase.eventPublisher.Publish(event)
}

func (*RequestPasswordResetUseCase) RequiredPermissions() []string {
	return make([]string, 0)
}

func NewRequestPasswordResetUseCase(userRepository user.UserRepository, userPasswordResetRepository user.UserPasswordResetRepository, eventPublisher events.EventPublisher, hasher hash.Hasher, logger *zap.Logger) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepository:              userRepository,
		userPasswordResetRepository: userPasswordResetRepository,
		hasher:                      hasher,
		eventPublisher:              eventPublisher,
		logger:                      logger,
	}
}
