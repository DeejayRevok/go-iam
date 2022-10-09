package thirdParty

type ThirdPartyTokensToEmailTransformer interface {
	Transform(tokens *ThirdPartyTokens) (string, error)
}
