package app

import (
	"go-uaa/src/domain/auth/thirdParty"
	"go-uaa/src/infrastructure/oauth2"
	"go-uaa/src/infrastructure/transformers"
	"os"
)

func BuildOauth2GoogleAuthURLBuilder() *oauth2.Oauth2GoogleAuthURLBuilder {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	state := os.Getenv("OAUTH_STATE")

	return oauth2.NewOauth2GoogleAuthURLBuilder(clientID, clientSecret, state)
}

func BuildOauth2GoogleTokensFetcher(tokensTransformer *transformers.Oauth2TokenToThirdPartyTokensTransformer) *oauth2.Oauth2GoogleTokensFetcher {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	return oauth2.NewOauth2GoogleTokensFetcher(clientID, clientSecret, tokensTransformer)
}

func BuildThirdPartyAuthStateChecker() *thirdParty.ThirdPartyAuthStateChecker {
	state := os.Getenv("OAUTH_STATE")

	return thirdParty.NewThirdPartyAuthStateChecker(state)
}
