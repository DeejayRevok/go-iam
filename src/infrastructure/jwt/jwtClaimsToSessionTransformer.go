package jwt

import (
	"errors"
	"go-uaa/src/domain/session"

	"github.com/golang-jwt/jwt"
)

type JWTClaimsToSessionTransformer struct{}

func (transformer *JWTClaimsToSessionTransformer) Transform(jwtClaims *jwt.MapClaims) (*session.Session, error) {
	claims := *jwtClaims

	userID, typeCheck := claims["userId"].(string)
	if !typeCheck {
		return nil, errors.New("user id token claim not valid")
	}
	exp, typeCheck := claims["exp"].(float64)
	if !typeCheck {
		return nil, errors.New("exp token claim not valid")
	}
	lastUsage, typeCheck := claims["lastUsage"].(float64)
	if !typeCheck {
		return nil, errors.New("last usage token claim not valid")
	}

	session := session.Session{
		UserID:    userID,
		Exp:       int64(exp),
		LastUsage: int64(lastUsage),
	}
	return &session, nil
}

func NewJWTClaimsToSessionTransformer() *JWTClaimsToSessionTransformer {
	transformer := JWTClaimsToSessionTransformer{}
	return &transformer
}
