package thirdParty

type ThirdPartyTokensFetcher interface {
	Fetch(code string, callbackURL string) (*ThirdPartyTokens, error)
}
