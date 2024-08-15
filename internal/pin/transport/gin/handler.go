package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
)

type handler struct {
	service utils.PinService
}

func NewHandler(service utils.PinService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/create", h.CreatePin())

}

func (h *handler) CreatePin() func(*gin.Context) {
	return func(c *gin.Context) {
		var pin models.PinCreation

		if err := c.ShouldBindJSON(&pin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := h.service.CreatePin(c.Request.Context(), &pin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(id))
	}
}

func (h *handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
