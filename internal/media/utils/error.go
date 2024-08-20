package utils

import "errors"

var (
	ErrNoFileReceived   = errors.New("No file received")
	ErrUnableToOpenFile = errors.New("Unable to open the file")
)
