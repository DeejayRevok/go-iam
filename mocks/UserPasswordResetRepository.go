// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	user "go-iam/src/domain/user"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserPasswordResetRepository is an autogenerated mock type for the UserPasswordResetRepository type
type UserPasswordResetRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, userPasswordReset
func (_m *UserPasswordResetRepository) Delete(ctx context.Context, userPasswordReset user.UserPasswordReset) error {
	ret := _m.Called(ctx, userPasswordReset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, user.UserPasswordReset) error); ok {
		r0 = rf(ctx, userPasswordReset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByUserID provides a mock function with given fields: ctx, userID
func (_m *UserPasswordResetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*user.UserPasswordReset, error) {
	ret := _m.Called(ctx, userID)

	var r0 *user.UserPasswordReset
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *user.UserPasswordReset); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.UserPasswordReset)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, userPasswordReset
func (_m *UserPasswordResetRepository) Save(ctx context.Context, userPasswordReset user.UserPasswordReset) error {
	ret := _m.Called(ctx, userPasswordReset)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, user.UserPasswordReset) error); ok {
		r0 = rf(ctx, userPasswordReset)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserPasswordResetRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserPasswordResetRepository creates a new instance of UserPasswordResetRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserPasswordResetRepository(t mockConstructorTestingTNewUserPasswordResetRepository) *UserPasswordResetRepository {
	mock := &UserPasswordResetRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
