package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"github.com/ssonit/aura_server/internal/auth/utils"
)

type handler struct {
	service utils.UserService
}

func NewHandler(service utils.UserService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/register", h.Register())
	group.POST("/login", h.Login())
}

func (h *handler) Login() func(*gin.Context) {
	return func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data, err := h.service.Login(c.Request.Context(), user.Email, user.Password)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

func (h *handler) Register() func(*gin.Context) {
	return func(c *gin.Context) {
		var user models.UserCreation

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data, err := h.service.Register(c.Request.Context(), &user)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusCreated, common.SimpleSuccessResponse(data))
	}
}