package sendPasswordResetToken

import (
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SendPasswordResetTokenUseCase struct {
	userRepository   user.UserRepository
	resetTokenSender user.PasswordResetTokenSender
	logger           *zap.Logger
}

func (useCase *SendPasswordResetTokenUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*SendPasswordResetTokenRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting sending password reset token %s to user %s", validatedRequest.ResetToken, validatedRequest.UserID))
	defer useCase.logger.Info(fmt.Sprintf("Finished sending password reset token %s to user %s", validatedRequest.ResetToken, validatedRequest.UserID))

	receiverID, err := uuid.Parse(validatedRequest.UserID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	receiver, err := useCase.userRepository.FindByID(receiverID)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	err = useCase.resetTokenSender.Send(validatedRequest.ResetToken, receiver)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.EmptyUseCaseResponse()
}

func (*SendPasswordResetTokenUseCase) RequiredPermissions() []string {
	return make([]string, 0)
}

func NewSendPasswordResetTokenUseCase(userRepo user.UserRepository, tokenSender user.PasswordResetTokenSender, logger *zap.Logger) *SendPasswordResetTokenUseCase {
	return &SendPasswordResetTokenUseCase{
		userRepository:   userRepo,
		resetTokenSender: tokenSender,
		logger:           logger,
	}
}
