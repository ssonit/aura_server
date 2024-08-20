package utils

import "errors"

var (
	ErrScanData      = errors.New("Failed to scan data")
	ErrUserIDIsBlank = errors.New("User id is blank")
)
