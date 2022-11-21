package resolvers

import (
	"context"
	"go-uaa/src/application/getAuthenticatedUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/user"
	"go-uaa/src/infrastructure/api"
	"go-uaa/src/infrastructure/graph/modelResolvers"
	"go-uaa/src/infrastructure/transformers"
	"net/http"
)

type MeResolver struct {
	getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase
	useCaseExecutor             *internals.AuthorizedUseCaseExecutor
	accessTokenFinder           *api.HTTPAccessTokenFinder
	userToResponseTransformer   *transformers.UserToResponseTransformer
}

func (resolver *MeResolver) Me(c context.Context) (*modelResolvers.UserResolver, error) {
	httpRequest := c.Value(RequestKey).(*http.Request)
	useCaseCtx := httpRequest.Context()
	accessToken, err := resolver.accessTokenFinder.Find(httpRequest)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	useCaseRequest := getAuthenticatedUser.GetAuthenticatedUserRequest{
		Token: accessToken,
	}
	useCaseResponse := resolver.useCaseExecutor.Execute(useCaseCtx, resolver.getAuthenticatedUserUseCase, &useCaseRequest, nil)
	if useCaseResponse.Err != nil {
		return nil, useCaseResponse.Err
	}
	userResponse := resolver.userToResponseTransformer.Transform(useCaseResponse.Content.(*user.User))
	return modelResolvers.NewUserResolver(*userResponse), nil
}

func NewMeResolver(getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, userToResponseTransformer *transformers.UserToResponseTransformer) *MeResolver {
	return &MeResolver{
		getAuthenticatedUserUseCase: getAuthenticatedUserUseCase,
		useCaseExecutor:             useCaseExecutor,
		accessTokenFinder:           accessTokenFinder,
		userToResponseTransformer:   userToResponseTransformer,
	}
}
