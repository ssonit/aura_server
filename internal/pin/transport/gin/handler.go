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
	group.GET("/", h.ListPinItem())
	group.GET("/:id", h.GetPinById())
	group.GET("/:id/board", h.ListBoardPinItem())

}

func (h *handler) ListBoardPinItem() func(*gin.Context) {
	return func(c *gin.Context) {

		boardId := c.Param("id")

		var data []models.BoardPinModel

		BoardId, err := primitive.ObjectIDFromHex(boardId)

		filter := models.BoardPinFilter{
			BoardId: BoardId,
		}

		data, err = h.service.ListBoardPinItem(c.Request.Context(), &filter)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, nil, filter, nil))
	}
}

func (h *handler) GetPinById() func(*gin.Context) {

	return func(c *gin.Context) {
		id := c.Param("id")

		data, err := h.service.GetPinById(c.Request.Context(), id)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func (h *handler) ListPinItem() func(*gin.Context) {
	return func(c *gin.Context) {
		var data []models.PinModel
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		paging.Process()

		var filter models.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		if err := filter.DecodeQuery(); err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.ListPinItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, filter, nil))

	}
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
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
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
