package authenticate

type AuthenticationRequest struct {
	Username     string
	Password     string
	Issuer       string
	GrantType    string
	RefreshToken string
}
