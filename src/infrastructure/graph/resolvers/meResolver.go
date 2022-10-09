package resolvers

import (
	"context"
	"go-uaa/src/application/getAuthenticatedUser"
	"go-uaa/src/domain/internals"
	"go-uaa/src/domain/session"
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
	sessionFinder               *api.HTTPSessionFinder
}

func (resolver *MeResolver) Me(c context.Context) (*modelResolvers.UserResolver, error) {
	httpRequest := c.Value(RequestKey).(*http.Request)
	useCaseCtx := httpRequest.Context()
	accessToken, err := resolver.accessTokenFinder.Find(httpRequest)
	if err != nil {
		return nil, err
	}
	var requestSession *session.Session
	if accessToken == nil {
		requestSession, err = resolver.sessionFinder.Find(httpRequest)
	}
	if err != nil {
		return nil, err
	}

	useCaseRequest := getAuthenticatedUser.GetAuthenticatedUserRequest{
		Token:   accessToken,
		Session: requestSession,
	}
	useCaseResponse := resolver.useCaseExecutor.Execute(useCaseCtx, resolver.getAuthenticatedUserUseCase, &useCaseRequest, nil, nil)
	if useCaseResponse.Err != nil {
		return nil, useCaseResponse.Err
	}
	userResponse := resolver.userToResponseTransformer.Transform(useCaseResponse.Content.(*user.User))
	return modelResolvers.NewUserResolver(*userResponse), nil
}

func NewMeResolver(getAuthenticatedUserUseCase *getAuthenticatedUser.GetAuthenticatedUserUseCase, useCaseExecutor *internals.AuthorizedUseCaseExecutor, accessTokenFinder *api.HTTPAccessTokenFinder, userToResponseTransformer *transformers.UserToResponseTransformer, sessionFinder *api.HTTPSessionFinder) *MeResolver {
	return &MeResolver{
		getAuthenticatedUserUseCase: getAuthenticatedUserUseCase,
		useCaseExecutor:             useCaseExecutor,
		accessTokenFinder:           accessTokenFinder,
		userToResponseTransformer:   userToResponseTransformer,
		sessionFinder:               sessionFinder,
	}
}
