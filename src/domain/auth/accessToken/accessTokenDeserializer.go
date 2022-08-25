package accessToken

type AccessTokenDeserializer interface {
	Deserialize(serializedToken string) (*AccessToken, error)
}
