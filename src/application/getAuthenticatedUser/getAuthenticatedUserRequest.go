package getAuthenticatedUser

import (
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/session"
)

type GetAuthenticatedUserRequest struct {
	Token   *accessToken.AccessToken
	Session *session.Session
}
