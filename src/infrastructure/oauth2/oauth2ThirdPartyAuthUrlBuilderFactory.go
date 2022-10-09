package oauth2

import (
	"fmt"
	"go-uaa/src/domain/auth/thirdParty"
)

const googleProvider = "google"

type Oauth2ThirdPartyAuthURLBuilderFactory struct {
	googleBuilder *Oauth2GoogleAuthURLBuilder
}

func (factory *Oauth2ThirdPartyAuthURLBuilderFactory) Create(provider string) (thirdParty.ThirdPartyAuthURLBuilder, error) {
	switch provider {
	case googleProvider:
		return factory.googleBuilder, nil
	default:
		return nil, fmt.Errorf("%s third party authentication provider not supported", provider)
	}
}

func NewOauth2ThirdPartyAuthURLBuilderFactory(googleBuilder *Oauth2GoogleAuthURLBuilder) *Oauth2ThirdPartyAuthURLBuilderFactory {
	return &Oauth2ThirdPartyAuthURLBuilderFactory{
		googleBuilder: googleBuilder,
	}
}
