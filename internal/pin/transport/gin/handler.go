package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/internal/pin/utils"
)

type handler struct {
	service utils.PinService
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(r *gin.Engine) {

	r.GET("/ping", h.Ping)
}

func (h *handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
