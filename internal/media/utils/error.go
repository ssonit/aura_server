package utils

import "errors"

var (
	ErrNoFileReceived     = errors.New("No file received")
	ErrUnableToOpenFile   = errors.New("Unable to open the file")
	ErrNotInserted        = errors.New("Not inserted")
	ErrCldNewFromParams   = errors.New("Failed to create new cloudinary instance")
	ErrCannotUploadCld    = errors.New("Failed to upload to cloudinary")
	ErrCannotCreateEntity = errors.New("Cannot create entity")
	ErrCannotGetEntity    = errors.New("Cannot get entity")
)
