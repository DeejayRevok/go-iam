package resolvers

import (
	"context"
	"go-iam/src/application/createUser"
	"go-iam/src/domain/internals"
	"go-iam/src/infrastructure/dto"
	"go-iam/src/infrastructure/graph/modelResolvers"
	"net/http"
)

type CreateUserResolver struct {
	createUserUseCase *createUser.CreateUserUseCase
	useCaseExecutor   *internals.UseCaseExecutor
}

func (resolver *CreateUserResolver) CreateUser(c context.Context, args *struct{ Input *dto.UserCreationRequestDTO }) (*modelResolvers.CreationResponse, error) {
	httpRequest := c.Value(RequestKey).(*http.Request)
	useCaseCtx := httpRequest.Context()
	createUserRequest := createUser.CreateUserRequest{
		Username: *args.Input.Username,
		Email:    *args.Input.Email,
		Password: *args.Input.Password,
	}
	useCaseResponse := resolver.useCaseExecutor.Execute(useCaseCtx, resolver.createUserUseCase, &createUserRequest, nil)
	if useCaseResponse.Err != nil {
		return modelResolvers.NewFailedCreationResponse(), useCaseResponse.Err
	}

	return modelResolvers.NewSuccessfulCreationResponse(), nil
}

func NewCreateUserResolver(createUserUseCase *createUser.CreateUserUseCase, useCaseExecutor *internals.UseCaseExecutor) *CreateUserResolver {
	return &CreateUserResolver{
		createUserUseCase: createUserUseCase,
		useCaseExecutor:   useCaseExecutor,
	}
}
