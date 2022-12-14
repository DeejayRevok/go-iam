package authenticationStrategy

import (
	"context"
	"go-iam/src/domain/user"
)

type AuthenticationStrategy interface {
	ValidateStrategyRequest(request *AuthenticationStrategyRequest) error
	Authenticate(ctx context.Context, request *AuthenticationStrategyRequest) (*user.User, error)
}
