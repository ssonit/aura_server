package utils

import "errors"

var (
	ErrScanData           = errors.New("Failed to scan data")
	ErrUserIDIsBlank      = errors.New("User id is blank")
	ErrNoDocuments        = errors.New("No documents found")
	ErrCannotCreateEntity = errors.New("Cannot create entity")
	ErrCannotGetEntity    = errors.New("Cannot get entity")
	ErrNotInserted        = errors.New("Failed to insert document")
	ErrDataIsZero         = errors.New("Data is zero")
)
