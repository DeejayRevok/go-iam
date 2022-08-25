package jwt

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/hex"
	"go-uaa/src/infrastructure/dto"
	"strconv"

	"github.com/golang-jwt/jwt"
)

type JWTRSAKeyToJWTKeyResponseTransformer struct{}

func (transformer *JWTRSAKeyToJWTKeyResponseTransformer) Transform(key []byte, keyUsage string) (*dto.JWTKeyDTO, error) {
	jwtRsaKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	keyDTO := dto.JWTKeyDTO{
		Kty: "RSA",
		E:   transformer.getKeyE(jwtRsaKey),
		Use: keyUsage,
		Kid: transformer.getKeyKid(key),
		Alg: transformer.getKeyAlg(),
		N:   transformer.getKeyN(jwtRsaKey),
	}
	return &keyDTO, nil
}

func (*JWTRSAKeyToJWTKeyResponseTransformer) getKeyE(jwtRsaKey *rsa.PublicKey) string {
	return strconv.Itoa(jwtRsaKey.E)
}

func (*JWTRSAKeyToJWTKeyResponseTransformer) getKeyAlg() string {
	return jwt.SigningMethodRS256.Alg()
}

func (*JWTRSAKeyToJWTKeyResponseTransformer) getKeyN(jwtRsaKey *rsa.PublicKey) string {
	return jwtRsaKey.N.String()
}

func (*JWTRSAKeyToJWTKeyResponseTransformer) getKeyKid(rawKey []byte) string {
	keyFingerprint := md5.Sum(rawKey)
	return hex.EncodeToString(keyFingerprint[:])
}

func NewJWTRSAKeyToJWTKeyResponseTransformer() *JWTRSAKeyToJWTKeyResponseTransformer {
	transformer := JWTRSAKeyToJWTKeyResponseTransformer{}
	return &transformer
}
