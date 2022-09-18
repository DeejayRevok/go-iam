package authenticationStrategy

import (
	"context"
	"errors"
	"go-uaa/src/domain/hash"
	"go-uaa/src/domain/user"
)

type PasswordAuthenticationStrategy struct {
	userRepository         user.UserRepository
	passwordHashComparator hash.HashComparator
}

func (strategy *PasswordAuthenticationStrategy) Authenticate(ctx context.Context, request *AuthenticationStrategyRequest) (*user.User, error) {
	user, err := strategy.userRepository.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, StrategyAuthenticationError{
			Username: request.Username,
			Strategy: "password",
			Message:  "User not found",
		}
	}
	err = strategy.passwordHashComparator.Compare(request.Password, user.Password)
	if err != nil {
		return nil, StrategyAuthenticationError{
			Username: request.Username,
			Strategy: "password",
			Message:  err.Error(),
		}
	}
	return user, nil
}

func (strategy *PasswordAuthenticationStrategy) ValidateStrategyRequest(request *AuthenticationStrategyRequest) error {
	if request.Username == "" {
		return errors.New("missing username for password authentication")
	}
	if request.Password == "" {
		return errors.New("missing password for password authentication")
	}
	return nil
}

func NewPasswordAuthenticationStrategy(userRepository user.UserRepository, passwordHashComparator hash.HashComparator) *PasswordAuthenticationStrategy {
	strategy := PasswordAuthenticationStrategy{
		userRepository:         userRepository,
		passwordHashComparator: passwordHashComparator,
	}
	return &strategy
}
