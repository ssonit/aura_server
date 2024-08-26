package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"github.com/ssonit/aura_server/internal/auth/utils"
)

var (
	jwtSecret = common.EnvConfig("JWT_SECRET", "secret")
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
		var user models.UserLogin

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.Login(c.Request.Context(), user.Email, user.Password)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		token, err := common.GenerateJWT([]byte(jwtSecret), data.ID.Hex())

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_TOKEN"))
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponseWithToken(data, token))

	}
}

func (h *handler) Register() func(*gin.Context) {
	return func(c *gin.Context) {
		var user models.UserCreation

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		data, err := h.service.Register(c.Request.Context(), &user)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		token, err := common.GenerateJWT([]byte(jwtSecret), data.ID.Hex())

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_TOKEN"))
			return
		}

		c.JSON(http.StatusCreated, common.SimpleSuccessResponseWithToken(data, token))
	}
}
