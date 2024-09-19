package utils

import (
	"github.com/ssonit/aura_server/common"
)

var (
	ErrScanData              = common.NewFullCustomError(500, "Failed to scan data", "SCAN_DATA_ERROR")
	ErrUserIDIsBlank         = common.NewFullCustomError(400, "User id is blank", "USER_ID_BLANK")
	ErrPinNotFound           = common.NewFullCustomError(404, "Pin not found", "PIN_NOT_FOUND")
	ErrCannotCreatePin       = common.NewFullCustomError(500, "Cannot create pin", "CANNOT_CREATE_PIN")
	ErrCannotGetEntity       = common.NewFullCustomError(500, "Cannot get entity", "CANNOT_GET_ENTITY")
	ErrDataIsZero            = common.NewFullCustomError(400, "Data is zero", "DATA_IS_ZERO")
	ErrFailedToFindEntity    = common.NewFullCustomError(500, "Failed to find entity", "FAILED_TO_FIND_ENTITY")
	ErrFailedToDecode        = common.NewFullCustomError(500, "Failed to decode", "FAILED_TO_DECODE")
	ErrCursorError           = common.NewFullCustomError(500, "Cursor error", "CURSOR_ERROR")
	ErrFailedToCount         = common.NewFullCustomError(500, "Failed to count", "FAILED_TO_COUNT")
	ErrCannotCreateBoardPin  = common.NewFullCustomError(500, "Cannot create board pin", "CANNOT_CREATE_BOARD_PIN")
	ErrUserNotPermitted      = common.NewFullCustomError(403, "User not permitted", "USER_NOT_PERMITTED")
	ErrCannotUpdatePin       = common.NewFullCustomError(500, "Cannot update pin", "CANNOT_UPDATE_PIN")
	ErrCannotDeleteBoardPin  = common.NewFullCustomError(500, "Cannot delete board pin", "CANNOT_DELETE_BOARD_PIN")
	ErrBoardNotFound         = common.NewFullCustomError(404, "Board not found", "BOARD_NOT_FOUND")
	ErrBoardPinAlreadyExists = common.NewFullCustomError(400, "Board pin already exists", "BOARD_PIN_ALREADY_EXISTS")
)
