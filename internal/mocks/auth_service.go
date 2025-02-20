// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/christhianjesus/crabi-challenge/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, email, password
func (_m *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	ret := _m.Called(ctx, email, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, email, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Signin provides a mock function with given fields: ctx, user
func (_m *AuthService) Signin(ctx context.Context, user *domain.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Signin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
