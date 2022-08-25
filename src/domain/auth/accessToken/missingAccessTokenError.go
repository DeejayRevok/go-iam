package accessToken

type MissingAccessTokenError struct{}

func (MissingAccessTokenError) Error() string {
	return "Missing access token"
}
