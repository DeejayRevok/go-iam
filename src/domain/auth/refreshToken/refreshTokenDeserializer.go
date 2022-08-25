package refreshToken

type RefreshTokenDeserializer interface {
	Deserialize(serializedToken string) (*RefreshToken, error)
}
