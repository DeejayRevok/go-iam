package auth

type MissingAuthorizationError struct{}

func (err MissingAuthorizationError) Error() string {
	return "Either access token or session should be provided"
}
