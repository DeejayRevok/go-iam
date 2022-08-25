package getAuthenticatedUser

import (
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"

	"go.uber.org/zap"
)

type GetAuthenticatedUserUseCase struct {
	userRepository user.UserRepository
	logger         *zap.Logger
}

func (useCase *GetAuthenticatedUserUseCase) Execute(request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*GetAuthenticatedUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info("Starting to get authenticated user data")
	defer useCase.logger.Info("Finished getting authenticated user data")

	user, err := useCase.userRepository.FindByUsername(validatedRequest.Token.Sub)
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	return internals.UseCaseResponse{
		Content: user,
		Err:     nil,
	}
}

func (*GetAuthenticatedUserUseCase) RequiredPermissions() []string {
	return []string{}
}

func NewGetAuthenticatedUserUseCase(userRepository user.UserRepository, logger *zap.Logger) *GetAuthenticatedUserUseCase {
	useCase := GetAuthenticatedUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
	return &useCase
}
