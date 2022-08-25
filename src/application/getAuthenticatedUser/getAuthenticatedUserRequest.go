package getAuthenticatedUser

import "go-uaa/src/domain/auth/accessToken"

type GetAuthenticatedUserRequest struct {
	Token accessToken.AccessToken
}
