package auth

type MissingAuthorizationError struct{}

func (err MissingAuthorizationError) Error() string {
	return "Access token should be provided"
}
