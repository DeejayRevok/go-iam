package jwt

import "go-uaa/src/infrastructure/dto"

type JWTKeySetBuilder struct {
	jwtSettings       *JWTSettings
	jwtKeyTransformer *JWTRSAKeyToJWTKeyResponseTransformer
}

func (builder *JWTKeySetBuilder) Build() (*dto.JWTKeySetResponseDTO, error) {
	jwtRsaPublicKey, err := builder.getJWTRSAPublicKey()
	if err != nil {
		return nil, err
	}

	keys := []dto.JWTKeyDTO{
		*jwtRsaPublicKey,
	}
	keySet := dto.JWTKeySetResponseDTO{
		Keys: keys,
	}
	return &keySet, nil
}

func (builder *JWTKeySetBuilder) getJWTRSAPublicKey() (*dto.JWTKeyDTO, error) {
	key, err := builder.jwtKeyTransformer.Transform(builder.jwtSettings.PublicKey, "sig")
	if err != nil {
		return nil, err
	}
	return key, nil
}

func NewJWTKeySetBuilder(jwtSettings *JWTSettings, jwtKetTransformer *JWTRSAKeyToJWTKeyResponseTransformer) *JWTKeySetBuilder {
	builder := JWTKeySetBuilder{
		jwtSettings:       jwtSettings,
		jwtKeyTransformer: jwtKetTransformer,
	}
	return &builder
}
