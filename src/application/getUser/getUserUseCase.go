package getUser

import (
	"context"
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
)

type GetUserUseCase struct {
	userRepository user.UserRepository
	logger         internals.Logger
}

func (useCase *GetUserUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*GetUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, fmt.Sprintf("Starting getting user %s", validatedRequest.UserId.String()))
	defer useCase.logger.Info(ctx, fmt.Sprintf("Finished getting user %s", validatedRequest.UserId.String()))

	user, err := useCase.userRepository.FindByID(ctx, validatedRequest.UserId)
	return internals.UseCaseResponse{
		Content: user,
		Err:     err,
	}
}

func (*GetUserUseCase) RequiredPermissions() []string {
	return []string{user.ReadUserPermission}
}

func NewGetUserUseCase(userRepository user.UserRepository, logger internals.Logger) *GetUserUseCase {
	useCase := GetUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
	return &useCase
}
