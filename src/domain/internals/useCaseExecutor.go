package internals

import (
	"context"
	"go-iam/src/domain/auth/accessToken"
)

type UseCaseExecutor struct{}

func (executor *UseCaseExecutor) Execute(ctx context.Context, useCase UseCase, useCaseRequest any, accessToken *accessToken.AccessToken) *UseCaseResponse {
	useCaseResponse := useCase.Execute(ctx, useCaseRequest)
	return &useCaseResponse
}

func NewAuthorizedUseCaseExecutor() *UseCaseExecutor {
	return &UseCaseExecutor{}
}
