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
	group.GET("/suggestions", h.ListSuggestions())
	group.GET("/soft-deleted", h.ListSoftDeletedPins())
	group.GET("/tags", h.ListTags())
	group.POST("/create", h.CreatePin())
	group.GET("/", h.ListPinItem())
	group.GET("/board-pin/:boardId/pins", h.ListBoardPinItem())
	group.GET("/board-pin/detail/:pinId", h.GetBoardPinItem())
	group.POST("/board-pin/save", h.SaveBoardPin())
	group.DELETE("/board-pin/:boardPinId/unsave/:pinId", h.UnSaveBoardPin())
	group.POST("/:id/like", h.LikePin())
	group.DELETE("/:id/unlike", h.UnlikePin())
	group.POST("/:id/comment", h.CreateComment())
	group.DELETE("/comment/:commentId", h.DeleteComment())
	group.GET("/:id/list/comments", h.ListComments())
	group.GET("/:id", h.GetPinById())
	group.PUT("/:id", h.UpdatePin())
	group.DELETE("/:id/soft-delete", h.SoftDeletePin())
	group.POST("/:id/restore", h.RestorePin())

}

func (h *handler) ListTags() func(*gin.Context) {
	return func(c *gin.Context) {

		var data []models.Tag
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		paging.Process()

		data, err := h.service.ListTags(c.Request.Context(), &paging)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, nil, nil))
	}
}

func (h *handler) ListSuggestions() func(*gin.Context) {
	return func(c *gin.Context) {
		keyword := c.Query("keyword")

		limit := 20

		data, err := h.service.ListSuggestions(c.Request.Context(), keyword, limit)

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

func (h *handler) ListSoftDeletedPins() func(*gin.Context) {
	return func(c *gin.Context) {

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.ListSoftDeletedPins(c.Request.Context(), userID.(string))

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

func (h *handler) RestorePin() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.RestorePin(c.Request.Context(), id, userID.(string))

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin restored successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}

}

func (h *handler) SoftDeletePin() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.SoftDeletePin(c.Request.Context(), id, userID.(string))

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin deleted successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func (h *handler) UnSaveBoardPin() func(*gin.Context) {
	return func(c *gin.Context) {

		boardPinId := c.Param("boardPinId")

		pinId := c.Param("pinId")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.UnSaveBoardPin(c.Request.Context(), &models.BoardPinUnSave{
			PinId:      pinId,
			UserId:     userID.(string),
			BoardPinId: boardPinId,
		})

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin deleted from board successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) DeleteComment() func(*gin.Context) {
	return func(c *gin.Context) {

		commentId := c.Param("commentId")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.DeleteComment(c.Request.Context(), commentId, userID.(string))

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Comment deleted successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) ListComments() func(*gin.Context) {
	return func(c *gin.Context) {
		pin_id := c.Param("id")

		var data []models.CommentModel
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		paging.Process()

		data, err := h.service.ListCommentsByPinId(c.Request.Context(), pin_id, &paging)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, nil, nil))
	}
}

func (h *handler) CreateComment() func(*gin.Context) {
	return func(c *gin.Context) {

		pin_id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var comment models.CommentCreation

		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		comment.PinId = pin_id
		comment.UserId = userID.(string)

		id, err := h.service.CreateComment(c.Request.Context(), &comment)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"id": id,
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func (h *handler) UnlikePin() func(*gin.Context) {
	return func(c *gin.Context) {

		pin_id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var like models.LikeDelete

		like.PinId = pin_id
		like.UserId = userID.(string)

		err := h.service.UnLikePin(c.Request.Context(), &like)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin unliked successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) LikePin() func(*gin.Context) {
	return func(c *gin.Context) {

		pin_id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var like models.LikeCreation

		like.PinId = pin_id
		like.UserId = userID.(string)

		err := h.service.LikePin(c.Request.Context(), &like)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin liked successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) SaveBoardPin() func(*gin.Context) {
	return func(c *gin.Context) {
		var boardPin models.BoardPinSave

		if err := c.ShouldBindJSON(&boardPin); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		boardPin.UserId, _ = userID.(string)

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.SaveBoardPin(c.Request.Context(), &boardPin)

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

func (h *handler) GetBoardPinItem() func(*gin.Context) {

	return func(c *gin.Context) {
		pinId := c.Param("pinId")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var filter models.BoardPinFilter
		var err error

		filter.PinId, err = primitive.ObjectIDFromHex(pinId)
		filter.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.GetBoardPinItem(c.Request.Context(), &filter)

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

func (h *handler) ListBoardPinItem() func(*gin.Context) {
	return func(c *gin.Context) {

		boardId := c.Param("boardId")

		var data []models.BoardPinModel
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		paging.Process()

		BoardId, err := primitive.ObjectIDFromHex(boardId)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		filter := models.BoardPinFilter{
			BoardId: BoardId,
		}

		data, err = h.service.ListBoardPinItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, filter, nil))
	}
}

func (h *handler) UpdatePin() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		var pin models.PinUpdate

		if err := c.ShouldBindJSON(&pin); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.UpdatePin(c.Request.Context(), id, &pin, userID.(string))

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Pin updated successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}

func (h *handler) GetPinById() func(*gin.Context) {

	return func(c *gin.Context) {
		id := c.Param("id")

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.GetPinById(c.Request.Context(), id, userID.(string))

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

func (h *handler) ListPinItem() func(*gin.Context) {
	return func(c *gin.Context) {
		var data []models.PinModel
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		paging.Process()

		var filter models.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		if err := filter.DecodeQuery(); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.ListPinItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, filter, nil))

	}
}

func (h *handler) CreatePin() func(*gin.Context) {
	return func(c *gin.Context) {
		var pin models.PinCreation

		if err := c.ShouldBindJSON(&pin); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		userID, exists := c.Get("userID")

		if !exists {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, utils.ErrUserIDIsBlank.Error(), "INVALID_REQUEST"))
			return
		}

		var err error

		pin.UserId, err = primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		id, err := h.service.CreatePin(c.Request.Context(), &pin)

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
