package authenticationStrategy

import "go-uaa/src/domain/auth/thirdParty"

type AuthenticationStrategyRequest struct {
	Username              string
	Password              string
	RefreshToken          string
	ThirdPartyAuthRequest *thirdParty.ThirdPartyAuthRequest
}
