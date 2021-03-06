// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"
	db "genesis/db"

	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/gofrs/uuid"
)

// AuthProvider is an autogenerated mock type for the AuthProvider type
type AuthProvider struct {
	mock.Mock
}

// GenerateJWT provides a mock function with given fields: ctx, user, userID, roleID, userAgent, expiration
func (_m *AuthProvider) GenerateJWT(ctx context.Context, user string, userID string, roleID string, userAgent string, expiration *time.Time) (string, error) {
	ret := _m.Called(ctx, user, userID, roleID, userAgent, expiration)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, *time.Time) string); ok {
		r0 = rf(ctx, user, userID, roleID, userAgent, expiration)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, *time.Time) error); ok {
		r1 = rf(ctx, user, userID, roleID, userAgent, expiration)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleFromContext provides a mock function with given fields: ctx, incTrackActions
func (_m *AuthProvider) RoleFromContext(ctx context.Context, incTrackActions bool) (*db.Role, error) {
	ret := _m.Called(ctx, incTrackActions)

	var r0 *db.Role
	if rf, ok := ret.Get(0).(func(context.Context, bool) *db.Role); ok {
		r0 = rf(ctx, incTrackActions)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, incTrackActions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleIDFromContext provides a mock function with given fields: ctx
func (_m *AuthProvider) RoleIDFromContext(ctx context.Context) (uuid.UUID, error) {
	ret := _m.Called(ctx)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context) uuid.UUID); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserFromContext provides a mock function with given fields: ctx
func (_m *AuthProvider) UserFromContext(ctx context.Context) (*db.User, error) {
	ret := _m.Called(ctx)

	var r0 *db.User
	if rf, ok := ret.Get(0).(func(context.Context) *db.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*db.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserIDFromContext provides a mock function with given fields: ctx
func (_m *AuthProvider) UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	ret := _m.Called(ctx)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context) uuid.UUID); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidatePassword provides a mock function with given fields: ctx, email, password
func (_m *AuthProvider) ValidatePassword(ctx context.Context, email string, password string) error {
	ret := _m.Called(ctx, email, password)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
