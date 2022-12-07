package refreshToken

import "go-iam/src/domain/user"

type RefreshTokenRequest struct {
	User *user.User
}
