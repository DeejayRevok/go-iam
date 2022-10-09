package oauth2

import (
	"context"
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/infrastructure/transformers"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Oauth2GoogleTokensFetcher struct {
	googleOAuthClientID     string
	googleOAuthClientSecret string
	tokenTransformer        *transformers.Oauth2TokenToThirdPartyTokensTransformer
}

func (fetcher *Oauth2GoogleTokensFetcher) Fetch(code string, callbackURL string) (*thirdParty.ThirdPartyTokens, error) {
	oauthConfing := &oauth2.Config{
		ClientID:     fetcher.googleOAuthClientID,
		ClientSecret: fetcher.googleOAuthClientSecret,
		RedirectURL:  callbackURL,
		Scopes:       []string{googleEmailScope},
		Endpoint:     google.Endpoint,
	}
	response, err := oauthConfing.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return fetcher.tokenTransformer.Transform(response), nil
}

func NewOauth2GoogleTokensFetcher(googleOAuthClientID string, googleOAuthClientSecret string, tokenTransformer *transformers.Oauth2TokenToThirdPartyTokensTransformer) *Oauth2GoogleTokensFetcher {
	return &Oauth2GoogleTokensFetcher{
		googleOAuthClientID:     googleOAuthClientID,
		googleOAuthClientSecret: googleOAuthClientSecret,
		tokenTransformer:        tokenTransformer,
	}
}
