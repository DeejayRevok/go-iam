package thirdParty

import "errors"

type ThirdPartyAuthStateChecker struct {
	internalState string
}

func (checker *ThirdPartyAuthStateChecker) Check(state string) error {
	if state != checker.internalState {
		return errors.New("auth state is not valid")
	}
	return nil
}

func NewThirdPartyAuthStateChecker(internalState string) *ThirdPartyAuthStateChecker {
	return &ThirdPartyAuthStateChecker{
		internalState: internalState,
	}
}
