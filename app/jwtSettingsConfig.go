package app

import (
	"go-iam/src/infrastructure/jwt"
	"io/ioutil"
	"os"
)

func LoadJWTSettings() *jwt.JWTSettings {
	privateKey, err := ioutil.ReadFile(os.Getenv("IAM_JWT_RSA_PRIVATE_KEY_FILE"))
	if err != nil {
		panic(err)
	}
	publicKey, err := ioutil.ReadFile(os.Getenv("IAM_JWT_RSA_PUBLIC_KEY_FILE"))
	if err != nil {
		panic(err)
	}

	return jwt.NewJWTSettings(privateKey, publicKey)
}
