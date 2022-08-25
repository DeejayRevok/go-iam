package transformers

import (
	"go-uaa/src/domain/auth/refreshToken"

	"github.com/golang-jwt/jwt"
)

type RefreshTokenToJWTClaimsTransformer struct{}

func (transformer *RefreshTokenToJWTClaimsTransformer) Transform(refreshToken *refreshToken.RefreshToken) *jwt.MapClaims {
	claims := make(jwt.MapClaims)

	claims["id"] = refreshToken.Id
	claims["sub"] = refreshToken.Sub
	claims["exp"] = refreshToken.Exp

	return &claims
}

func NewRefreshTokenToJWTClaimsTransformer() *RefreshTokenToJWTClaimsTransformer {
	transformer := RefreshTokenToJWTClaimsTransformer{}
	return &transformer
}
