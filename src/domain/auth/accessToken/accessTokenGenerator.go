package accessToken

import (
	"time"
)

type AccessTokenGenerator struct{}

func (generator *AccessTokenGenerator) Generate(request *AccessTokenRequest) AccessToken {
	iat := time.Now()
	exp := iat.Add(time.Hour * time.Duration(AccessTokenDefaultExpirationHours))
	return AccessToken{
		Iss:   request.Issuer,
		Sub:   request.User.Email,
		Iat:   iat.Unix(),
		Exp:   exp.Unix(),
		Scope: "",
	}
}

func NewAccessTokenGenerator() *AccessTokenGenerator {
	generator := AccessTokenGenerator{}
	return &generator
}
