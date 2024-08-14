package utils

import "errors"

var (
	ErrScanData              = errors.New("failed to scan data")
	ErrInvalidTodoItemStatus = errors.New("invalid todo item status")
	ErrTitleIsBlank          = errors.New("title is blank")
	ErrTodoItemIsDeleted     = errors.New("todo item is deleted")
)
