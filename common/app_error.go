package common

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewFullErrorResponse(statusCode int, rootErr error, message, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    rootErr,
		Message:    message,
		Log:        log,
		Key:        key,
	}
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "database error", err.Error(), "DB_ERROR")
}
func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "invalid request", err.Error(), "INVALID_REQUEST")
}

func ErrInternalServer(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "internal server error", err.Error(), "INTERNAL_SERVER_ERROR")
}

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot create %s", entity), "CANNOT_CREATE_ENTITY")
}

func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot get %s", entity), "CANNOT_GET_ENTITY")
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot update %s", entity), "CANNOT_UPDATE_ENTITY")
}

func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot delete %s", entity), "CANNOT_DELETE_ENTITY")
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s not found", entity), "ENTITY_NOT_FOUND")
}

func ErrNoPermission(err error) *AppError {
	return NewCustomError(err, "No permission", "NO_PERMISSION")
}

func NewErrorResponse(rootErr error, message, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    rootErr,
		Log:        log,
		Message:    message,
		Key:        key,
	}
}

func NewUnauthorizedResponse(rootErr error, message, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    rootErr,
		Message:    message,
		Key:        key,
	}
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}

	return e.RootErr
}

func (e *AppError) Error() string {
	return e.RootError().Error()
}

func NewCustomError(root error, message, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, message, root.Error(), key)
	}
	return NewErrorResponse(errors.New(message), message, message, key)
}

var RecordNotFound = errors.New("record not found")
