// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	thirdParty "go-iam/src/domain/auth/thirdParty"

	mock "github.com/stretchr/testify/mock"
)

// ThirdPartyTokensFetcher is an autogenerated mock type for the ThirdPartyTokensFetcher type
type ThirdPartyTokensFetcher struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: code, callbackURL
func (_m *ThirdPartyTokensFetcher) Fetch(code string, callbackURL string) (*thirdParty.ThirdPartyTokens, error) {
	ret := _m.Called(code, callbackURL)

	var r0 *thirdParty.ThirdPartyTokens
	if rf, ok := ret.Get(0).(func(string, string) *thirdParty.ThirdPartyTokens); ok {
		r0 = rf(code, callbackURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*thirdParty.ThirdPartyTokens)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(code, callbackURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewThirdPartyTokensFetcher interface {
	mock.TestingT
	Cleanup(func())
}

// NewThirdPartyTokensFetcher creates a new instance of ThirdPartyTokensFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewThirdPartyTokensFetcher(t mockConstructorTestingTNewThirdPartyTokensFetcher) *ThirdPartyTokensFetcher {
	mock := &ThirdPartyTokensFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
