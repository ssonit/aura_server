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
	group.GET("/:id", h.GetMedia())
}

func (h *handler) GetMedia() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		data, err := h.service.GetMedia(c.Request.Context(), id)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func (h *handler) UploadImage() func(*gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.UploadImage(c.Request.Context(), file)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		res := map[string]interface{}{
			"id": id,
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(res))
	}
}
