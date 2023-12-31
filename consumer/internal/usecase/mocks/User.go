// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	user "go-nats-pub-sub-restapi-postgresql/consumer/internal/entity"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// AccrualBalanceUser provides a mock function with given fields: ctx, userDTO
func (_m *User) AccrualBalanceUser(ctx context.Context, userDTO *user.User) error {
	ret := _m.Called(ctx, userDTO)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) error); ok {
		r0 = rf(ctx, userDTO)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CalcNewBalance provides a mock function with given fields: ctx, userDTO, balance
func (_m *User) CalcNewBalance(ctx context.Context, userDTO *user.User, balance float32) float32 {
	ret := _m.Called(ctx, userDTO, balance)

	var r0 float32
	if rf, ok := ret.Get(0).(func(context.Context, *user.User, float32) float32); ok {
		r0 = rf(ctx, userDTO, balance)
	} else {
		r0 = ret.Get(0).(float32)
	}

	return r0
}

// CheckNegativeBalance provides a mock function with given fields: ctx, userDTO
func (_m *User) CheckNegativeBalance(ctx context.Context, userDTO *user.User) error {
	ret := _m.Called(ctx, userDTO)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) error); ok {
		r0 = rf(ctx, userDTO)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, userDTO
func (_m *User) CreateUser(ctx context.Context, userDTO *user.User) error {
	ret := _m.Called(ctx, userDTO)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) error); ok {
		r0 = rf(ctx, userDTO)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBalance provides a mock function with given fields: ctx, id
func (_m *User) GetBalance(ctx context.Context, id uint64) (user.User, error) {
	ret := _m.Called(ctx, id)

	var r0 user.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (user.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) user.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(user.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUser creates a new instance of User. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUser(t interface {
	mock.TestingT
	Cleanup(func())
}) *User {
	mock := &User{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
