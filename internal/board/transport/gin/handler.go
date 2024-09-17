package gin

import (
	"net/http"
	"strconv"

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
	group.GET("/:id", h.GetBoardItem())
	group.POST("/create", h.CreateBoard())
	group.GET("/", h.ListBoardItem())
	group.DELETE("/:id/soft-delete", h.SoftDeleteBoard())
	group.POST("/:id/restore", h.RestoreBoard())
	group.PUT("/:id", h.UpdateBoardItem())
	group.GET("/deleted", h.ListDeletedBoards())

}

func (h *handler) ListDeletedBoards() func(*gin.Context) {
	return func(c *gin.Context) {
		var data []models.BoardModel

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.ListDeletedBoards(c.Request.Context(), userID.(string))

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

func (h *handler) RestoreBoard() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrBoardIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.RestoreBoard(c.Request.Context(), id)

		if err != nil {

			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Board soft deleted successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func (h *handler) SoftDeleteBoard() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrBoardIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.SoftDeleteBoard(c.Request.Context(), id)

		if err != nil {

			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Board soft deleted successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) UpdateBoardItem() func(*gin.Context) {
	return func(c *gin.Context) {
		var board models.BoardUpdate

		if err := c.ShouldBindJSON(&board); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		id := c.Param("id")

		if id == "" {

			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrBoardIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.UpdateBoardItem(c.Request.Context(), id, &board)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Board updated successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) GetBoardItem() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrBoardIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		boardID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		board, err := h.service.GetBoardItem(c.Request.Context(), boardID)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(board))
	}
}

func (h *handler) CreateBoard() func(*gin.Context) {

	return func(c *gin.Context) {
		var board models.BoardCreation

		if err := c.ShouldBindJSON(&board); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var err error

		board.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.CreateBoard(c.Request.Context(), &board)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
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

		var privateValue bool

		userId := c.Query("userId")
		isPrivate := c.Query("isPrivate")

		privateValue, err = strconv.ParseBool(isPrivate)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		if privateValue {
			filter.IsPrivate = privateValue
		}

		if userId == "" {
			userIDValue, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
				return
			}

			userId = userIDValue.(string)

		}

		filter.UserId, err = primitive.ObjectIDFromHex(userId)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err = h.service.ListBoardItem(c.Request.Context(), &filter)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, nil, filter, nil))

	}
}
