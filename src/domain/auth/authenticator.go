package auth

import (
	"context"
	"fmt"
	"go-iam/src/domain/auth/authenticationStrategy"
	"go-iam/src/domain/user"
	"strings"
)

type Authenticator struct {
	passwordAuthenticationStrategy     *authenticationStrategy.PasswordAuthenticationStrategy
	refreshTokenAuthenticationStrategy *authenticationStrategy.RefreshTokenAuthenticationStrategy
	thirdPartyAuthenticationStrategy   *authenticationStrategy.ThirdPartyAuthenticationStrategy
}

func (authenticator *Authenticator) Authenticate(ctx context.Context, grantType string, request *authenticationStrategy.AuthenticationStrategyRequest) (*user.User, error) {
	strategy, err := authenticator.getStrategy(grantType)
	if err != nil {
		return nil, err
	}
	if err = strategy.ValidateStrategyRequest(request); err != nil {
		return nil, err
	}
	user, err := strategy.Authenticate(ctx, request)
	if err == nil {
		return user, nil
	}
	if _, isinstance := err.(authenticationStrategy.StrategyAuthenticationError); isinstance {
		return nil, AuthenticationError{
			Email: request.Email,
		}
	}
	return nil, err
}

func (authenticator *Authenticator) getStrategy(grantType string) (authenticationStrategy.AuthenticationStrategy, error) {
	switch grantType = strings.ToLower(grantType); grantType {
	case "password":
		return authenticator.passwordAuthenticationStrategy, nil
	case "refresh_token":
		return authenticator.refreshTokenAuthenticationStrategy, nil
	case "third_party":
		return authenticator.thirdPartyAuthenticationStrategy, nil
	default:
		return nil, fmt.Errorf("grant type %s is not supported", grantType)
	}
}

func NewAuthenticator(passwordAuthenticationStrategy *authenticationStrategy.PasswordAuthenticationStrategy, refreshTokenAuthenticationStrategy *authenticationStrategy.RefreshTokenAuthenticationStrategy, thirdPartyAuthenticationStrategy *authenticationStrategy.ThirdPartyAuthenticationStrategy) *Authenticator {
	authenticator := Authenticator{
		passwordAuthenticationStrategy:     passwordAuthenticationStrategy,
		refreshTokenAuthenticationStrategy: refreshTokenAuthenticationStrategy,
		thirdPartyAuthenticationStrategy:   thirdPartyAuthenticationStrategy,
	}
	return &authenticator
}
