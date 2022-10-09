package jwt

import (
	"errors"
	"go-uaa/src/domain/auth/thirdParty"

	"github.com/golang-jwt/jwt"
)

type JWTThirdPartyTokensToEmailTransformer struct{}

func (*JWTThirdPartyTokensToEmailTransformer) Transform(tokens *thirdParty.ThirdPartyTokens) (string, error) {
	jwtToken, _ := jwt.Parse(tokens.IDToken, nil)
	if err := jwtToken.Claims.Valid(); err != nil {
		return "", err
	}
	tokenClaims, typeCheck := jwtToken.Claims.(jwt.MapClaims)
	if !typeCheck {
		return "", errors.New("token claims not valid")
	}
	return tokenClaims["email"].(string), nil
}

func NewJWTThirdPartyTokensToEmailTransformer() *JWTThirdPartyTokensToEmailTransformer {
	return &JWTThirdPartyTokensToEmailTransformer{}
}
