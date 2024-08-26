package utils

import (
	"github.com/ssonit/aura_server/common"
)

var (
	ErrEmailAlreadyExists   = common.NewFullCustomError(409, "Email already exists", "EMAIL_EXISTS")
	ErrFailedToHashPassword = common.NewFullCustomError(500, "Failed to hash password", "HASHING_ERROR")
	ErrCannotCreateEntity   = common.NewFullCustomError(500, "Failed to create user", "DATABASE_ERROR")
	ErrUserNotFound         = common.NewFullCustomError(404, "User not found", "USER_NOT_FOUND")
	ErrEmailOrPassInvalid   = common.NewFullCustomError(400, "Email or password invalid", "INVALID_CREDENTIALS")
)
