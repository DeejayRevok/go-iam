package jwt

import (
	"github.com/golang-jwt/jwt"
)

type JWTTokenGenerator struct {
	settings *JWTSettings
}

func (encoder *JWTTokenGenerator) Generate(tokenClaims *jwt.MapClaims) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(encoder.settings.PrivateKey)
	if err != nil {
		return "", err
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func NewJWTTokenGenerator(settings *JWTSettings) *JWTTokenGenerator {
	generator := JWTTokenGenerator{
		settings: settings,
	}
	return &generator
}
