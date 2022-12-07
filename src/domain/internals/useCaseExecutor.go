package internals

import (
	"context"
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/user"
)

type UseCaseExecutor struct {
	userRepository user.UserRepository
}

func (executor *UseCaseExecutor) Execute(ctx context.Context, useCase UseCase, useCaseRequest any, accessToken *accessToken.AccessToken) *UseCaseResponse {
	useCaseResponse := useCase.Execute(ctx, useCaseRequest)
	return &useCaseResponse
}

func NewAuthorizedUseCaseExecutor() *UseCaseExecutor {
	return &UseCaseExecutor{}
}
