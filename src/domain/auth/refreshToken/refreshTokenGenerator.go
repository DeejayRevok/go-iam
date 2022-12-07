package refreshToken

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenGenerator struct{}

func (*RefreshTokenGenerator) Generate(request *RefreshTokenRequest) RefreshToken {
	iat := time.Now()
	exp := iat.Add(time.Hour * time.Duration(RefreshTokenDefaultExpirationHours))
	id := uuid.New()
	return RefreshToken{
		Sub: request.User.Email,
		Exp: exp.Unix(),
		Id:  id.String(),
	}
}

func NewRefreshTokenGenerator() *RefreshTokenGenerator {
	generator := RefreshTokenGenerator{}
	return &generator
}
