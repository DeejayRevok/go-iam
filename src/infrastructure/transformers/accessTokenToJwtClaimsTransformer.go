package transformers

import (
	"go-iam/src/domain/auth/accessToken"

	"github.com/golang-jwt/jwt"
)

type AccessTokenToJWTClaimsTransformer struct{}

func (transformer *AccessTokenToJWTClaimsTransformer) Transform(acessToken *accessToken.AccessToken) *jwt.MapClaims {
	claims := make(jwt.MapClaims)

	claims["iss"] = acessToken.Iss
	claims["sub"] = acessToken.Sub
	claims["exp"] = acessToken.Exp
	claims["iat"] = acessToken.Iat
	claims["scope"] = acessToken.Scope

	return &claims
}

func NewAccessTokenToJWTClaimsTransformer() *AccessTokenToJWTClaimsTransformer {
	transformer := AccessTokenToJWTClaimsTransformer{}
	return &transformer
}
