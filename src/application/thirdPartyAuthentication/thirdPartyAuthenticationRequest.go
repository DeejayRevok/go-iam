package thirdPartyAuthentication

type ThirdPartyAuthenticationRequest struct {
	State        string
	Code         string
	AuthProvider string
	CallbackURL  string
}
