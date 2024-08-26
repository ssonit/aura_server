package utils

import (
	"github.com/ssonit/aura_server/common"
)

var (
	ErrNoFileReceived     = common.NewFullCustomError(400, "No file received", "NO_FILE_RECEIVED")
	ErrUnableToOpenFile   = common.NewFullCustomError(500, "Unable to open the file", "FILE_OPEN_ERROR")
	ErrCldNewFromParams   = common.NewFullCustomError(500, "Failed to create new Cloudinary instance", "CLOUDINARY_INSTANCE_ERROR")
	ErrCannotUploadCld    = common.NewFullCustomError(500, "Failed to upload to Cloudinary", "CLOUDINARY_UPLOAD_ERROR")
	ErrCannotCreateEntity = common.NewFullCustomError(500, "Cannot create entity", "CANNOT_CREATE_ENTITY")
	ErrCannotGetEntity    = common.NewFullCustomError(500, "Cannot get entity", "CANNOT_GET_ENTITY")
)
