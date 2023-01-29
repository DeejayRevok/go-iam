package app

import (
	"go-iam/src/domain/auth/thirdParty"
	"go-iam/src/infrastructure/oauth2"
	"go-iam/src/infrastructure/transformers"
	"os"
)

func BuildOauth2GoogleAuthURLBuilder() *oauth2.Oauth2GoogleAuthURLBuilder {
	clientID := os.Getenv("IAM_GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("IAM_GOOGLE_OAUTH_CLIENT_SECRET")
	state := os.Getenv("IAM_OAUTH_STATE")

	return oauth2.NewOauth2GoogleAuthURLBuilder(clientID, clientSecret, state)
}

func BuildOauth2GoogleTokensFetcher(tokensTransformer *transformers.Oauth2TokenToThirdPartyTokensTransformer) *oauth2.Oauth2GoogleTokensFetcher {
	clientID := os.Getenv("IAM_GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("IAM_GOOGLE_OAUTH_CLIENT_SECRET")

	return oauth2.NewOauth2GoogleTokensFetcher(clientID, clientSecret, tokensTransformer)
}

func BuildThirdPartyAuthStateChecker() *thirdParty.ThirdPartyAuthStateChecker {
	state := os.Getenv("IAM_OAUTH_STATE")

	return thirdParty.NewThirdPartyAuthStateChecker(state)
}
