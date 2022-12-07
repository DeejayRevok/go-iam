package authenticate

type AuthenticationRequest struct {
	Email                  string
	Password               string
	Issuer                 string
	GrantType              string
	RefreshToken           string
	ThirdPartyState        string
	ThirdPartyCode         string
	ThirdPartyAuthProvider string
	ThirdPartyCallbackURL  string
}
