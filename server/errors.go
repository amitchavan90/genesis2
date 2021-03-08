package genesis

import (
	"errors"
)

// ErrCreateLimitReached when the limit has been reached
var ErrCreateLimitReached = errors.New("limit reached")

// ErrPasswordShort when the entered password is too short
var ErrPasswordShort = errors.New("password too shirt")

// ErrPasswordCommon when the entered password is too common
var ErrPasswordCommon = errors.New("password too common")

// ErrEmailInvalid when the entered email is invalid
var ErrEmailInvalid = errors.New("invalid email")

// ErrMobileNotSet when the current user has no mobile set
var ErrMobileNotSet = errors.New("mobile not set")

// ErrTokenInvalid when the token is invalid
var ErrTokenInvalid = errors.New("token invalid")

// ErrTokenExpired when the token is expired
var ErrTokenExpired = errors.New("token expired")

// ErrArchived for trying to update something that is archived
var ErrArchived = errors.New("object is archived")
