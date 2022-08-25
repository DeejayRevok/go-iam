package auth

import "fmt"

type AuthenticationError struct {
	Username string
}

func (err AuthenticationError) Error() string {
	return fmt.Sprintf("Error authenticating %s", err.Username)
}
