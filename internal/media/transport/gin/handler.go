package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/media/utils"
	"github.com/ssonit/aura_server/middleware"
)

type handler struct {
	service utils.MediaService
}

func NewHandler(service utils.MediaService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(group *gin.RouterGroup) {
	group.Use(middleware.AuthMiddleware())
	group.POST("/upload-image", h.UploadImage())
}

func (h *handler) UploadImage() func(*gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(utils.ErrNoFileReceived, utils.ErrNoFileReceived.Error(), "NO_FILE_RECEIVED"))
			return
		}

		id, err := h.service.UploadImage(c.Request.Context(), file)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "UPLOAD_FAILED"))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(id))
	}
}
