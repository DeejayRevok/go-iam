package internals

import (
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/user"
)

type AuthorizedUseCaseExecutor struct {
	userRepository user.UserRepository
}

func (executor *AuthorizedUseCaseExecutor) Execute(useCase UseCase, useCaseRequest any, accessToken *accessToken.AccessToken) *UseCaseResponse {
	requiredPermissions := useCase.RequiredPermissions()
	if len(requiredPermissions) > 0 {
		if err := executor.checkPermissions(useCase, accessToken, requiredPermissions); err != nil {
			useCaseResponse := UseCaseResponse{
				Err: err,
			}
			return &useCaseResponse
		}
	}

	useCaseResponse := useCase.Execute(useCaseRequest)
	return &useCaseResponse
}

func (executor *AuthorizedUseCaseExecutor) checkPermissions(useCase UseCase, token *accessToken.AccessToken, permissions []string) error {
	if token == nil {
		return accessToken.MissingAccessTokenError{}
	}
	user, err := executor.userRepository.FindByUsername(token.Sub)
	if err != nil {
		return err
	}

	if user.Superuser == true {
		return nil
	}

	for _, permissionName := range permissions {
		if user.HasPermission(permissionName) == false {
			return UseCaseAuthorizationError{
				Username:   user.Username,
				Permission: permissionName,
			}
		}
	}

	return nil
}

func NewAuthorizedUseCaseExecutor(userRepository user.UserRepository) *AuthorizedUseCaseExecutor {
	return &AuthorizedUseCaseExecutor{
		userRepository: userRepository,
	}
}
