package oauth2

import (
	"fmt"
	"go-iam/src/domain/auth/thirdParty"
)

type Oauth2ThirdPartyTokensFetcherFactory struct {
	googleFetcher *Oauth2GoogleTokensFetcher
}

func (factory *Oauth2ThirdPartyTokensFetcherFactory) Create(provider string) (thirdParty.ThirdPartyTokensFetcher, error) {
	switch provider {
	case googleProvider:
		return factory.googleFetcher, nil
	default:
		return nil, fmt.Errorf("%s third party authentication provider not supported", provider)
	}
}

func NewOauth2ThirdPartyTokensFetcherFactory(googleFetcher *Oauth2GoogleTokensFetcher) *Oauth2ThirdPartyTokensFetcherFactory {
	return &Oauth2ThirdPartyTokensFetcherFactory{
		googleFetcher: googleFetcher,
	}
}
