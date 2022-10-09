package oauth2

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleEmailScope = "email"

type Oauth2GoogleAuthURLBuilder struct {
	googleOAuthClientID     string
	googleOAuthClientSecret string
	state                   string
}

func (builder *Oauth2GoogleAuthURLBuilder) Build(callbackURL string) string {
	oauthConfing := &oauth2.Config{
		ClientID:     builder.googleOAuthClientID,
		ClientSecret: builder.googleOAuthClientSecret,
		RedirectURL:  callbackURL,
		Scopes:       []string{googleEmailScope},
		Endpoint:     google.Endpoint,
	}
	return oauthConfing.AuthCodeURL(builder.state)
}

func NewOauth2GoogleAuthURLBuilder(googleOAuthClientID string, googleOAuthClientSecret string, state string) *Oauth2GoogleAuthURLBuilder {
	return &Oauth2GoogleAuthURLBuilder{
		googleOAuthClientID:     googleOAuthClientID,
		googleOAuthClientSecret: googleOAuthClientSecret,
		state:                   state,
	}
}
