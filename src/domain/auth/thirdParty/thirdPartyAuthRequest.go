package thirdParty

type ThirdPartyAuthRequest struct {
	State        string
	Code         string
	AuthProvider string
	CallbackURL  string
}
