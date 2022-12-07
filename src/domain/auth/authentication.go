package auth

import (
	"go-iam/src/domain/auth/accessToken"
	"go-iam/src/domain/auth/refreshToken"
)

const DefaultTokenType = "bearer"

type Authentication struct {
	AccessToken  *accessToken.AccessToken
	RefreshToken *refreshToken.RefreshToken
	TokenType    string
	ExpiresIn    int
}
