package authenticationStrategy

import "fmt"

type StrategyAuthenticationError struct {
	Email    string
	Strategy string
	Message  string
}

func (err StrategyAuthenticationError) Error() string {
	return fmt.Sprintf("Error authenticating %s with strategy %s: %s", err.Email, err.Strategy, err.Message)
}
