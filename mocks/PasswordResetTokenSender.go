// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	user "go-iam/src/domain/user"

	mock "github.com/stretchr/testify/mock"
)

// PasswordResetTokenSender is an autogenerated mock type for the PasswordResetTokenSender type
type PasswordResetTokenSender struct {
	mock.Mock
}

// Send provides a mock function with given fields: resetToken, receiver
func (_m *PasswordResetTokenSender) Send(resetToken string, receiver *user.User) error {
	ret := _m.Called(resetToken, receiver)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *user.User) error); ok {
		r0 = rf(resetToken, receiver)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPasswordResetTokenSender interface {
	mock.TestingT
	Cleanup(func())
}

// NewPasswordResetTokenSender creates a new instance of PasswordResetTokenSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPasswordResetTokenSender(t mockConstructorTestingTNewPasswordResetTokenSender) *PasswordResetTokenSender {
	mock := &PasswordResetTokenSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
