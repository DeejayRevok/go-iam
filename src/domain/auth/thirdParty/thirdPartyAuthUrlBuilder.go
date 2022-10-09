package thirdParty

type ThirdPartyAuthURLBuilder interface {
	Build(callbackURL string) string
}
