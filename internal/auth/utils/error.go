package utils

import "errors"

var (
	ErrEmailAlreadyExists     = errors.New("Email already exists")
	ErrFailedToHashPassword   = errors.New("Failed to hash password")
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
	ErrNoDocuments            = errors.New("No documents found")
	ErrNotInserted            = errors.New("Failed to insert document")
	ErrCannotCreateEntity     = errors.New("Cannot create entity")
	ErrCannotGetEntity        = errors.New("Cannot get entity")
)
