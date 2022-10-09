package authenticate

import (
	"go-uaa/src/domain/auth"
	"go-uaa/src/domain/session"
)

type AuthenticationResponse struct {
	Authentication *auth.Authentication
	Session        *session.Session
}
