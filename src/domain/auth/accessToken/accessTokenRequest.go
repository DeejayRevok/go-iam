package accessToken

import "go-iam/src/domain/user"

type AccessTokenRequest struct {
	User   *user.User
	Issuer string
}
