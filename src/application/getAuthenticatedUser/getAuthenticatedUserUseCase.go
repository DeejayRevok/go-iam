package getAuthenticatedUser

import (
	"context"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
)

type GetAuthenticatedUserUseCase struct {
	userRepository user.UserRepository
	logger         internals.Logger
}

func (useCase *GetAuthenticatedUserUseCase) Execute(ctx context.Context, request any) internals.UseCaseResponse {
	validatedRequest, errResponse := internals.ValidateUseCaseRequest[*GetAuthenticatedUserRequest](request)
	if errResponse != nil {
		return *errResponse
	}

	useCase.logger.Info(ctx, "Starting to get authenticated user data")
	defer useCase.logger.Info(ctx, "Finished getting authenticated user data")

	user, err := useCase.userRepository.FindByUsername(ctx, validatedRequest.Token.Sub)
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

func NewGetAuthenticatedUserUseCase(userRepository user.UserRepository, logger internals.Logger) *GetAuthenticatedUserUseCase {
	useCase := GetAuthenticatedUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
	return &useCase
}
