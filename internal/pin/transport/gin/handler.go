package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"github.com/ssonit/aura_server/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// private
	group.Use(middleware.AuthMiddleware())
	group.POST("/create", h.CreatePin())

}

func (h *handler) CreatePin() func(*gin.Context) {
	return func(c *gin.Context) {
		var pin models.PinCreation

		if err := c.ShouldBindJSON(&pin); err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewCustomError(utils.ErrUserIDIsBlank, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var err error

		pin.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(common.InvalidObjectID, common.InvalidObjectID.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.CreatePin(c.Request.Context(), &pin)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
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
