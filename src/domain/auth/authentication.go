package auth

import (
	"go-uaa/src/domain/auth/accessToken"
	"go-uaa/src/domain/auth/refreshToken"
)

const DefaultTokenType = "bearer"

type Authentication struct {
	AccessToken  *accessToken.AccessToken
	RefreshToken *refreshToken.RefreshToken
	TokenType    string
	ExpiresIn    int
}
