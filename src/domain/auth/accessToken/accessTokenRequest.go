package accessToken

import "go-uaa/src/domain/user"

type AccessTokenRequest struct {
	User   *user.User
	Issuer string
}
