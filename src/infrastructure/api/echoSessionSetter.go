package api

import (
	"go-uaa/src/domain/session"
	"go-uaa/src/infrastructure/jwt"
	"go-uaa/src/infrastructure/transformers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoSessionSetter struct {
	sessionTransformer *transformers.SessionToJWTClaimsTransformer
	tokenGenerator     *jwt.JWTTokenGenerator
}

func (setter *EchoSessionSetter) Set(c echo.Context, session *session.Session) (echo.Context, error) {
	sessionClaims := setter.sessionTransformer.Transform(session)
	sessionJWTToken, err := setter.tokenGenerator.Generate(sessionClaims)
	if err != nil {
		return nil, err
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = sessionJWTToken
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c, nil
}

func NewEchoSessionSetter(sessionTransformer *transformers.SessionToJWTClaimsTransformer, tokenGenerator *jwt.JWTTokenGenerator) *EchoSessionSetter {
	return &EchoSessionSetter{
		sessionTransformer: sessionTransformer,
		tokenGenerator:     tokenGenerator,
	}
}
