package utils

import "errors"

var (
	ErrEmailAlreadyExists     = errors.New("Email already exists")
	ErrFailedToHashPassword   = errors.New("Failed to hash password")
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
)
