package transformers

import (
	"go-uaa/src/domain/auth"
	"go-uaa/src/infrastructure/dto"
	"go-uaa/src/infrastructure/jwt"
)

type AuthenticationToResponseTransformer struct {
	accessTokenTransformer  *AccessTokenToJWTClaimsTransformer
	refreshTokenTransformer *RefreshTokenToJWTClaimsTransformer
	tokenGenerator          *jwt.JWTTokenGenerator
}

func (transformer *AuthenticationToResponseTransformer) Transform(authentication *auth.Authentication) (*dto.AuthenticationResponseDTO, error) {
	accessTokenClaims := transformer.accessTokenTransformer.Transform(authentication.AccessToken)
	jwtAccessToken, err := transformer.tokenGenerator.Generate(accessTokenClaims)
	if err != nil {
		return nil, err
	}
	refreshTokenClaims := transformer.refreshTokenTransformer.Transform(authentication.RefreshToken)
	jwtRefreshToken, err := transformer.tokenGenerator.Generate(refreshTokenClaims)
	if err != nil {
		return nil, err
	}
	response := dto.AuthenticationResponseDTO{
		AccessToken:  jwtAccessToken,
		RefreshToken: jwtRefreshToken,
		TokenType:    authentication.TokenType,
		ExpiresIn:    authentication.ExpiresIn,
		Scope:        authentication.AccessToken.Scope,
	}
	return &response, nil
}

func NewAuthenticationToResponseTransformer(accessTokenTransformer *AccessTokenToJWTClaimsTransformer, refreshTokenTransformer *RefreshTokenToJWTClaimsTransformer, tokenGenerator *jwt.JWTTokenGenerator) *AuthenticationToResponseTransformer {
	transformer := AuthenticationToResponseTransformer{
		accessTokenTransformer:  accessTokenTransformer,
		refreshTokenTransformer: refreshTokenTransformer,
		tokenGenerator:          tokenGenerator,
	}
	return &transformer
}
