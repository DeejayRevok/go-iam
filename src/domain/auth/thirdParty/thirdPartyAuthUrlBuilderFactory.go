package thirdParty

type ThirdPartyAuthURLBuilderFactory interface {
	Create(provider string) (ThirdPartyAuthURLBuilder, error)
}
