package refreshToken

import "go-uaa/src/domain/user"

type RefreshTokenRequest struct {
	User *user.User
}
