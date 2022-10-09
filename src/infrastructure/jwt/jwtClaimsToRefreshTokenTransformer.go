package jwt

import (
	"errors"
	"go-uaa/src/domain/auth/refreshToken"

	"github.com/golang-jwt/jwt"
)

type JWTClaimsToRefreshTokenTransformer struct{}

func (*JWTClaimsToRefreshTokenTransformer) Transform(jwtClaims *jwt.MapClaims) (*refreshToken.RefreshToken, error) {
	claims := *jwtClaims

	sub, typeCheck := claims["sub"].(string)
	if !typeCheck {
		return nil, errors.New("sub token claim not valid")
	}
	exp, typeCheck := claims["exp"].(float64)
	if !typeCheck {
		return nil, errors.New("exp token claim not valid")
	}
	id, typeCheck := claims["id"].(string)
	if !typeCheck {
		return nil, errors.New("id token claim not valid")
	}

	token := refreshToken.RefreshToken{
		Sub: sub,
		Exp: int64(exp),
		Id:  id,
	}
	return &token, nil
}

func NewJWTClaimsToRefreshTokenTransformer() *JWTClaimsToRefreshTokenTransformer {
	transformer := JWTClaimsToRefreshTokenTransformer{}
	return &transformer
}
