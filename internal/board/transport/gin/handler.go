package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/board/models"
	"github.com/ssonit/aura_server/internal/board/utils"
	"github.com/ssonit/aura_server/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	service utils.BoardService
}

func NewHandler(service utils.BoardService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(group *gin.RouterGroup) {

	// private
	group.Use(middleware.AuthMiddleware())
	group.POST("/create", h.CreateBoard())
	group.GET("/", h.ListBoardItem())

}

func (h *handler) CreateBoard() func(*gin.Context) {

	return func(c *gin.Context) {
		var board models.BoardCreation

		if err := c.ShouldBindJSON(&board); err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewCustomError(utils.ErrUserIDIsBlank, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var err error

		board.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.CreateBoard(c.Request.Context(), &board)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(id))
	}
}

func (h *handler) ListBoardItem() func(*gin.Context) {
	return func(c *gin.Context) {
		var data []models.BoardModel

		var filter models.Filter

		var err error

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewCustomError(utils.ErrUserIDIsBlank, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		filter.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		data, err = h.service.ListBoardItem(c.Request.Context(), &filter)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewCustomError(err, err.Error(), "INVALID_REQUEST"))
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, nil, filter, nil))

	}
}
