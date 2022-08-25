package authenticationStrategy

import "fmt"

type StrategyAuthenticationError struct {
	Username string
	Strategy string
	Message  string
}

func (err StrategyAuthenticationError) Error() string {
	return fmt.Sprintf("Error authenticating %s with strategy %s: %s", err.Username, err.Strategy, err.Message)
}
