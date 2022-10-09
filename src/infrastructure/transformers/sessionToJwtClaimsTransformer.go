package transformers

import (
	"go-uaa/src/domain/session"

	"github.com/golang-jwt/jwt"
)

type SessionToJWTClaimsTransformer struct{}

func (transformer *SessionToJWTClaimsTransformer) Transform(session *session.Session) *jwt.MapClaims {
	claims := make(jwt.MapClaims)

	claims["userId"] = session.UserID
	claims["exp"] = session.Exp
	claims["lastUsage"] = session.LastUsage

	return &claims
}

func NewSessionToJWTClaimsTransformer() *SessionToJWTClaimsTransformer {
	transformer := SessionToJWTClaimsTransformer{}
	return &transformer
}
