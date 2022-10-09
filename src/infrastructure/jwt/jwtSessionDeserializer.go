package jwt

import (
	"errors"
	"go-uaa/src/domain/session"

	"github.com/golang-jwt/jwt"
)

type JWTSessionDeserializer struct {
	settings         *JWTSettings
	tokenTransformer *JWTClaimsToSessionTransformer
}

func (deserializer *JWTSessionDeserializer) Deserialize(serializedToken string) (*session.Session, error) {
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
		return nil, errors.New("session token claims not valid")
	}
	return deserializer.tokenTransformer.Transform(&tokenClaims)
}

func NewJWTSessionDeserializer(settings *JWTSettings, tokenTransformer *JWTClaimsToSessionTransformer) *JWTSessionDeserializer {
	deserializer := JWTSessionDeserializer{
		settings:         settings,
		tokenTransformer: tokenTransformer,
	}
	return &deserializer
}
