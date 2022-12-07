package auth

import "fmt"

type AuthenticationError struct {
	Email string
}

func (err AuthenticationError) Error() string {
	return fmt.Sprintf("Error authenticating %s", err.Email)
}
