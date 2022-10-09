package jwt

import (
	"errors"
	"go-uaa/src/domain/auth/accessToken"

	"github.com/golang-jwt/jwt"
)

type JWTClaimsToAccessTokenTransformer struct{}

func (transformer *JWTClaimsToAccessTokenTransformer) Transform(jwtClaims *jwt.MapClaims) (*accessToken.AccessToken, error) {
	claims := *jwtClaims

	iss, typeCheck := claims["iss"].(string)
	if !typeCheck {
		return nil, errors.New("iss token claim not valid")
	}
	sub, typeCheck := claims["sub"].(string)
	if !typeCheck {
		return nil, errors.New("sub token claim not valid")
	}
	exp, typeCheck := claims["exp"].(float64)
	if !typeCheck {
		return nil, errors.New("exp token claim not valid")
	}
	iat, typeCheck := claims["iat"].(float64)
	if !typeCheck {
		return nil, errors.New("iat token claim not valid")
	}
	scope, typeCheck := claims["scope"].(string)
	if !typeCheck {
		return nil, errors.New("scope token claim not valid")
	}

	token := accessToken.AccessToken{
		Iss:   iss,
		Sub:   sub,
		Exp:   int64(exp),
		Iat:   int64(iat),
		Scope: scope,
	}
	return &token, nil
}

func NewJWTClaimsToAccessTokenTransformer() *JWTClaimsToAccessTokenTransformer {
	transformer := JWTClaimsToAccessTokenTransformer{}
	return &transformer
}
