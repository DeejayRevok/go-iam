package transformers

import (
	"go-uaa/src/domain/auth/thirdParty"

	"golang.org/x/oauth2"
)

type Oauth2TokenToThirdPartyTokensTransformer struct{}

func (transformer *Oauth2TokenToThirdPartyTokensTransformer) Transform(token *oauth2.Token) *thirdParty.ThirdPartyTokens {
	return &thirdParty.ThirdPartyTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		IDToken:      token.Extra("id_token").(string),
	}
}

func NewOauth2TokenToThirdPartyTokensTransformer() *Oauth2TokenToThirdPartyTokensTransformer {
	return &Oauth2TokenToThirdPartyTokensTransformer{}
}
