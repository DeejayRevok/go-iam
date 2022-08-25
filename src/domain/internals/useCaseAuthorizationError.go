package internals

import "fmt"

type UseCaseAuthorizationError struct {
	Username   string
	Permission string
}

func (err UseCaseAuthorizationError) Error() string {
	return fmt.Sprintf("User %s has not authorization for %s", err.Username, err.Permission)
}
