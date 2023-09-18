package admin

import "errors"

var (
	// ErrExpiredToken is returned when the given token is
	// expired.
	ErrExpiredToken = errors.New("expired token")

	// ErrInvalidPassword is returned when the given password
	// is invalid.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidToken is returned when the given token is
	// invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidEmail is returned when the given email is
	// invalid.
	ErrInvalidEmail = errors.New("invalid email")
)
