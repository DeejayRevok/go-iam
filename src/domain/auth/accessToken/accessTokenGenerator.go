package accessToken

import (
	"go-uaa/src/domain/user"
	"strings"
	"time"
)

type AccessTokenGenerator struct{}

func (generator *AccessTokenGenerator) Generate(request *AccessTokenRequest) AccessToken {
	iat := time.Now()
	exp := iat.Add(time.Hour * time.Duration(AccessTokenDefaultExpirationHours))
	return AccessToken{
		Iss:   request.Issuer,
		Sub:   request.User.Username,
		Iat:   iat.Unix(),
		Exp:   exp.Unix(),
		Scope: generator.getScopes(request.User),
	}
}

func (generator *AccessTokenGenerator) getScopes(user *user.User) string {
	scopes := make([]string, 0)
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			scopes = append(scopes, permission.Name)
		}
	}
	return strings.Join(scopes[:], " ")
}

func NewAccessTokenGenerator() *AccessTokenGenerator {
	generator := AccessTokenGenerator{}
	return &generator
}
