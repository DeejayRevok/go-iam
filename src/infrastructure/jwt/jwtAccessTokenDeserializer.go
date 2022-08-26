package jwt

import (
	"errors"
	"go-uaa/src/domain/auth/accessToken"

	"github.com/golang-jwt/jwt"
)

type JWTAccessTokenDeserializer struct {
	settings         *JWTSettings
	tokenTransformer *JWTClaimsToAccessTokenTransformer
}

func (deserializer *JWTAccessTokenDeserializer) Deserialize(serializedToken string) (*accessToken.AccessToken, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(deserializer.settings.PublicKey)
	if err != nil {
		return nil, err
	}
	jwtToken, err := jwt.Parse(serializedToken, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if err := jwtToken.Claims.Valid(); err != nil {
		return nil, err
	}
	tokenClaims, typeCheck := jwtToken.Claims.(jwt.MapClaims)
	if !typeCheck {
		return nil, errors.New("Token claims not valid")
	}
	return deserializer.tokenTransformer.Transform(&tokenClaims)
}

func NewJWTAccessTokenDeserializer(settings *JWTSettings, tokenTransformer *JWTClaimsToAccessTokenTransformer) *JWTAccessTokenDeserializer {
	deserializer := JWTAccessTokenDeserializer{
		settings:         settings,
		tokenTransformer: tokenTransformer,
	}
	return &deserializer
}
