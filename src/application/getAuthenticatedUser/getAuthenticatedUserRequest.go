package getAuthenticatedUser

import (
	"go-iam/src/domain/auth/accessToken"
)

type GetAuthenticatedUserRequest struct {
	Token *accessToken.AccessToken
}
