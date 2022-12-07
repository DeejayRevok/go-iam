package sendPasswordResetToken

import (
	"context"
	"fmt"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"

	"github.com/google/uuid"
)

type SendPasswordResetTokenUseCase struct {
	userRepository   user.UserRepository
	resetTokenSender user.PasswordResetTokenSender
	logger           internals.Logger
}

func (useCase *SendPasswordResetTokenUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*SendPasswordResetTokenRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting sending password reset token %s to user %s", validatedRequest.ResetToken, validatedRequest.UserID))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished sending password reset token %s to user %s", validatedRequest.ResetToken, validatedRequest.UserID))

	receiverID, err := uuid.Parse(validatedRequest.UserID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	receiver, err := useCase.userRepository.FindByID(ctx, receiverID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.resetTokenSender.Send(validatedRequest.ResetToken, receiver)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func NewSendPasswordResetTokenUseCase(userRepo user.UserRepository, tokenSender user.PasswordResetTokenSender, logger internals.Logger) *SendPasswordResetTokenUseCase {
	return &SendPasswordResetTokenUseCase{
		userRepository:   userRepo,
		resetTokenSender: tokenSender,
		logger:           logger,
	}
}
