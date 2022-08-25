package jwt

type JWTSettings struct {
	PrivateKey []byte
	PublicKey  []byte
}

func NewJWTSettings(privateKey []byte, publickey []byte) *JWTSettings {
	settings := JWTSettings{
		PrivateKey: privateKey,
		PublicKey:  publickey,
	}
	return &settings
}
