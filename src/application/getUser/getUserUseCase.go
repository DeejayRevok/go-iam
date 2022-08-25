package getUser

import (
	"fmt"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"

	"go.uber.org/zap"
)

type GetUserUseCase struct {
	userRepository user.UserRepository
	logger         *zap.Logger
}

func (useCase *GetUserUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*GetUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(fmt.Sprintf("Starting getting user %s", validatedRequest.UserId.String()))
	defer useCase.logger.Info(fmt.Sprintf("Finished getting user %s", validatedRequest.UserId.String()))

	user, err := useCase.userRepository.FindByID(validatedRequest.UserId)
	return internals.UseCaseResponse{
		Content: user,
		Err:     err,
	}
}

func (*GetUserUseCase) RequiredPermissions() []string {
	return []string{user.ReadUserPermission}
}

func NewGetUserUseCase(userRepository user.UserRepository, logger *zap.Logger) *GetUserUseCase {
	useCase := GetUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
	return &useCase
}
