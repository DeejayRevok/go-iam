package getAuthenticatedUser

import (
	"context"
	"errors"
	"go-uaa/src/domain/auth/accessToken"
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

	var user *user.User
	var err error
	if validatedRequest.Token != nil {
		user, err = useCase.getUserFromAccessToken(ctx, validatedRequest.Token)
	}
	if err != nil {
		return internals.ErrorUseCaseResponse(err)
	}
	if user == nil {
		return internals.ErrorUseCaseResponse(errors.New("missing authenticated user"))
	}
	return internals.UseCaseResponse{
		Content: user,
		Err:     nil,
	}
}

func (useCase *GetAuthenticatedUserUseCase) getUserFromAccessToken(ctx context.Context, token *accessToken.AccessToken) (*user.User, error) {
	return useCase.userRepository.FindByUsername(ctx, token.Sub)
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
