package authenticationStrategy

import "go-iam/src/domain/auth/thirdParty"

type AuthenticationStrategyRequest struct {
	Email                 string
	Password              string
	RefreshToken          string
	ThirdPartyAuthRequest *thirdParty.ThirdPartyAuthRequest
}
