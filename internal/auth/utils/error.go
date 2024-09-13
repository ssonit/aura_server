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
	ErrCannotLogout         = common.NewFullCustomError(500, "Failed to logout", "LOGOUT_ERROR")
	ErrUserIDIsBlank        = common.NewFullCustomError(400, "User id is blank", "USER_ID_BLANK")
	ErrRefreshTokenNotFound = common.NewFullCustomError(404, "Refresh token not found", "REFRESH_TOKEN")
	ErrFailedToFindEntity   = common.NewFullCustomError(500, "Failed to find entity", "FAILED_TO_FIND_ENTITY")
	ErrFailedToDecode       = common.NewFullCustomError(500, "Failed to decode", "FAILED_TO_DECODE")
	ErrCursorError          = common.NewFullCustomError(500, "Cursor error", "CURSOR_ERROR")
	ErrCannotUpdateUser     = common.NewFullCustomError(500, "Failed to update user", "UPDATE_USER_ERROR")
)
