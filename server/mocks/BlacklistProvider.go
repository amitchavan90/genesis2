// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// BlacklistProvider is an autogenerated mock type for the BlacklistProvider type
type BlacklistProvider struct {
	mock.Mock
}

// BlacklistAll provides a mock function with given fields: userID
func (_m *BlacklistProvider) BlacklistAll(userID string) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BlacklistOne provides a mock function with given fields: tokenID
func (_m *BlacklistProvider) BlacklistOne(tokenID string) error {
	ret := _m.Called(tokenID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tokenID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CleanIssuedTokens provides a mock function with given fields:
func (_m *BlacklistProvider) CleanIssuedTokens() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnList provides a mock function with given fields: tokenID
func (_m *BlacklistProvider) OnList(tokenID string) bool {
	ret := _m.Called(tokenID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(tokenID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RefreshBlacklist provides a mock function with given fields:
func (_m *BlacklistProvider) RefreshBlacklist() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartTicker provides a mock function with given fields: _a0
func (_m *BlacklistProvider) StartTicker(_a0 context.Context) {
	_m.Called(_a0)
}
