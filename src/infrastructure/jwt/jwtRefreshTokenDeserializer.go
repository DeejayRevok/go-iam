package jwt

import (
	"errors"
	"go-iam/src/domain/auth/refreshToken"

	"github.com/golang-jwt/jwt"
)

type JWTRefreshTokenDeserializer struct {
	settings         *JWTSettings
	tokenTransformer *JWTClaimsToRefreshTokenTransformer
}

func (deserializer *JWTRefreshTokenDeserializer) Deserialize(serializedToken string) (*refreshToken.RefreshToken, error) {
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
		return nil, errors.New("token claims not valid")
	}
	return deserializer.tokenTransformer.Transform(&tokenClaims)
}

func NewJWTRefreshTokenDeserializer(settings *JWTSettings, tokenTransformer *JWTClaimsToRefreshTokenTransformer) *JWTRefreshTokenDeserializer {
	deserializer := JWTRefreshTokenDeserializer{
		settings:         settings,
		tokenTransformer: tokenTransformer,
	}
	return &deserializer
}
