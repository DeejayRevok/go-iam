// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	thirdParty "go-uaa/src/domain/auth/thirdParty"

	mock "github.com/stretchr/testify/mock"
)

// ThirdPartyAuthURLBuilderFactory is an autogenerated mock type for the ThirdPartyAuthURLBuilderFactory type
type ThirdPartyAuthURLBuilderFactory struct {
	mock.Mock
}

// Create provides a mock function with given fields: provider
func (_m *ThirdPartyAuthURLBuilderFactory) Create(provider string) (thirdParty.ThirdPartyAuthURLBuilder, error) {
	ret := _m.Called(provider)

	var r0 thirdParty.ThirdPartyAuthURLBuilder
	if rf, ok := ret.Get(0).(func(string) thirdParty.ThirdPartyAuthURLBuilder); ok {
		r0 = rf(provider)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(thirdParty.ThirdPartyAuthURLBuilder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(provider)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewThirdPartyAuthURLBuilderFactory interface {
	mock.TestingT
	Cleanup(func())
}

// NewThirdPartyAuthURLBuilderFactory creates a new instance of ThirdPartyAuthURLBuilderFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewThirdPartyAuthURLBuilderFactory(t mockConstructorTestingTNewThirdPartyAuthURLBuilderFactory) *ThirdPartyAuthURLBuilderFactory {
	mock := &ThirdPartyAuthURLBuilderFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
