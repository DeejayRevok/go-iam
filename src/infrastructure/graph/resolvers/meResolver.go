package resolvers

import (
	"context"
	"go-iam/src/application/getAuthenticatedUser"
	"go-iam/src/domain/internals"
	"go-iam/src/domain/user"
	"go-iam/src/infrastructure/api"
	"go-iam/src/infrastructure/graph/modelResolvers"
	"go-iam/src/infrastructure/transformers"
	"net/http"
)

type MeResolver struct {
	getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase
	useCaseExecutor             *internals.UseCaseExecutor
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

func NewMeResolver(getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase, useCaseExecutor *internals.UseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, userToResponseTransformer *transformers.UserToResponseTransformer) *MeResolver {
	return &MeResolver{
		getAuthenticatedUserUseCase: getAuthenticatedUserUseCase,
		useCaseExecutor:             useCaseExecutor,
		accessTokenFinder:           accessTokenFinder,
		userToResponseTransformer:   userToResponseTransformer,
	}
}
