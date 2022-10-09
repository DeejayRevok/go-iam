package thirdParty

type ThirdPartyTokensFetcherFactory interface {
	Create(provider string) (ThirdPartyTokensFetcher, error)
}
