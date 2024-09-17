package utils

import (
	"github.com/ssonit/aura_server/common"
)

var (
	ErrUserIDIsBlank        = common.NewFullCustomError(400, "User id is blank", "USER_ID_BLANK")
	ErrBoardNotFound        = common.NewFullCustomError(404, "Board not found", "BOARD_NOT_FOUND")
	ErrCannotCreateEntity   = common.NewFullCustomError(500, "Cannot create entity", "CANNOT_CREATE_ENTITY")
	ErrCannotGetEntity      = common.NewFullCustomError(500, "Cannot get entity", "CANNOT_GET_ENTITY")
	ErrFailedToFindEntity   = common.NewFullCustomError(500, "Failed to find entity", "FAILED_TO_FIND_ENTITY")
	ErrFailedToDecode       = common.NewFullCustomError(500, "Failed to decode", "FAILED_TO_DECODE")
	ErrBoardIDIsBlank       = common.NewFullCustomError(400, "Board id is blank", "BOARD_ID_BLANK")
	ErrFailedToDecodeObjID  = common.NewFullCustomError(500, "Failed to decode object id", "FAILED_TO_DECODE_OBJ_ID")
	ErrFailedToUpdateEntity = common.NewFullCustomError(500, "Failed to update entity", "FAILED_TO_UPDATE_ENTITY")
)
