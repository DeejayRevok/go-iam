package authenticationStrategy

import (
	"context"
	"errors"
	"go-iam/src/domain/hash"
	"go-iam/src/domain/user"
)

type PasswordAuthenticationStrategy struct {
	userRepository         user.UserRepository
	passwordHashComparator hash.HashComparator
}

func (strategy *PasswordAuthenticationStrategy) Authenticate(ctx context.Context, request *AuthenticationStrategyRequest) (*user.User, error) {
	user, err := strategy.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, StrategyAuthenticationError{
			Email:    request.Email,
			Strategy: "password",
			Message:  "User not found",
		}
	}
	err = strategy.passwordHashComparator.Compare(request.Password, user.Password)
	if err != nil {
		return nil, StrategyAuthenticationError{
			Email:    request.Email,
			Strategy: "password",
			Message:  err.Error(),
		}
	}
	return user, nil
}

func (strategy *PasswordAuthenticationStrategy) ValidateStrategyRequest(request *AuthenticationStrategyRequest) error {
	if request.Email == "" {
		return errors.New("missing email for password authentication")
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
