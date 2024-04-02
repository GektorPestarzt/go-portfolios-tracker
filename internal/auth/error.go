package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exist")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
