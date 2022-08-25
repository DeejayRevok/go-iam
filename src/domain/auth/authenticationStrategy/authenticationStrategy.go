package authenticationStrategy

import (
	"go-uaa/src/domain/user"
)

type AuthenticationStrategy interface {
	ValidateStrategyRequest(request *AuthenticationStrategyRequest) error
	Authenticate(request *AuthenticationStrategyRequest) (*user.User, error)
}
