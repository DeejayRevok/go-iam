package resolvers

import (
	"context"
	"go-uaa/src/application/createUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/graph/modelResolvers"
	"net/http"
)

type CreateUserResolver struct {
	createUserUseCase *createUser.CreateUserUseCase
	useCaseExecutor   *internals.AuthorizedUseCaseExecutor
}

func (resolver *CreateUserResolver) CreateUser(c context.Context, args *struct{ Input *dto.UserCreationRequestDTO }) (*modelResolvers.CreationResponse, error) {
	httpRequest := c.Value(RequestKey).(*http.Request)
	useCaseCtx := httpRequest.Context()
	createUserRequest := createUser.CreateUserRequest{
		Username: *args.Input.Username,
		Email:    *args.Input.Email,
		Password: *args.Input.Password,
		Roles:    resolver.parseRoles(*args.Input.Roles),
	}
	useCaseResponse := resolver.useCaseExecutor.Execute(useCaseCtx, resolver.createUserUseCase, &createUserRequest, nil)
	if useCaseResponse.Err != nil {
		return modelResolvers.NewFailedCreationResponse(), useCaseResponse.Err
	}

	return modelResolvers.NewSuccessfulCreationResponse(), nil
}

func (resolver *CreateUserResolver) parseRoles(roles []*string) []string {
	parsedRoles := make([]string, 0)
	for _, role := range roles {
		parsedRoles = append(parsedRoles, *role)
	}
	return parsedRoles
}

func NewCreateUserResolver(createUserUseCase *createUser.CreateUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor) *CreateUserResolver {
	return &CreateUserResolver{
		createUserUseCase: createUserUseCase,
		useCaseExecutor:   useCaseExecutor,
	}
}
