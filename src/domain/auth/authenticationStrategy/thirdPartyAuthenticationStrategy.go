package authenticationStrategy

import (
	"context"
	"errors"
	"fmt"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/domain/user"
)

type ThirdPartyAuthenticationStrategy struct {
	authStateChecker         *thirdParty.ThirdPartyAuthStateChecker
	tokensFetcherFactory     thirdParty.ThirdPartyTokensFetcherFactory
	tokensToEmailTransformer thirdParty.ThirdPartyTokensToEmailTransformer
	userRepository           user.UserRepository
}

func (strategy *ThirdPartyAuthenticationStrategy) Authenticate(ctx context.Context, request *AuthenticationStrategyRequest) (*user.User, error) {
	if err := strategy.authStateChecker.Check(request.ThirdPartyAuthRequest.State); err != nil {
		return nil, err
	}
	tokensFetcher, err := strategy.tokensFetcherFactory.Create(request.ThirdPartyAuthRequest.AuthProvider)
	if err != nil {
		return nil, err
	}

	tokens, err := tokensFetcher.Fetch(request.ThirdPartyAuthRequest.Code, request.ThirdPartyAuthRequest.CallbackURL)
	if err != nil {
		return nil, err
	}

	email, err := strategy.tokensToEmailTransformer.Transform(tokens)
	if err != nil {
		return nil, err
	}

	user, err := strategy.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with email %s not registered", email)
	}
	return user, nil
}

func (strategy *ThirdPartyAuthenticationStrategy) ValidateStrategyRequest(request *AuthenticationStrategyRequest) error {
	if request.ThirdPartyAuthRequest == nil || request.ThirdPartyAuthRequest.AuthProvider == "" || request.ThirdPartyAuthRequest.State == "" || request.ThirdPartyAuthRequest.Code == "" || request.ThirdPartyAuthRequest.CallbackURL == "" {
		return errors.New("missing third party authentication request params")
	}
	return nil
}

func NewThirdPartyAuthenticationStrategy(authStateChecker *thirdParty.ThirdPartyAuthStateChecker, tokensFetcherFactory thirdParty.ThirdPartyTokensFetcherFactory, tokensToEmailTransformer thirdParty.ThirdPartyTokensToEmailTransformer, userRepository user.UserRepository) *ThirdPartyAuthenticationStrategy {
	return &ThirdPartyAuthenticationStrategy{
		authStateChecker:         authStateChecker,
		tokensFetcherFactory:     tokensFetcherFactory,
		tokensToEmailTransformer: tokensToEmailTransformer,
		userRepository:           userRepository,
	}
}
