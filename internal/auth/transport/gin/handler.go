package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"github.com/ssonit/aura_server/internal/auth/utils"
	"github.com/ssonit/aura_server/middleware"
)

var (
	jwtSecret        = common.EnvConfig("JWT_SECRET", "secret")
	jwtRefreshSecret = common.EnvConfig("JWT_REFRESH_SECRET", "secret")
	jwtSecretExp     = common.EnvConfig("JWT_SECRET_EXP", "30")
	jwtRefreshExp    = common.EnvConfig("JWT_REFRESH_EXP", "24")
	expSecretTime    = 365 * 24 * 60 // 1 year
	expRefreshTime   = 24 * 30 * 12  // 1 year
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
	group.POST("/refresh-token", h.RefreshToken())

	group.Use(middleware.AuthMiddleware())
	group.GET("/me", h.Me())
	group.GET("/:id", h.GetUser())
	group.PUT("/:id", h.UpdateUser())
	group.POST("/logout", h.Logout())
}

func (h *handler) UpdateUser() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		var user models.UserUpdate

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.UpdateUser(c.Request.Context(), id, &user)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		result := map[string]interface{}{
			"message": "Updated user successfully",
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func (h *handler) GetUser() func(*gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		data, err := h.service.GetUser(c.Request.Context(), id)

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

func (h *handler) Me() func(*gin.Context) {
	return func(c *gin.Context) {

	}
}

func (h *handler) RefreshToken() func(*gin.Context) {
	return func(c *gin.Context) {
		var refreshToken models.RefreshTokenSelection

		if err := c.ShouldBindJSON(&refreshToken); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.Logout(c.Request.Context(), refreshToken.Token)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		// decode refresh token
		claims, err := common.DecodedToken(refreshToken.Token, []byte(jwtRefreshSecret))

		userID := claims["userID"].(string)

		expSecret := time.Now().Add(time.Minute * time.Duration(expSecretTime)).Unix()
		expRefresh := time.Now().Add(time.Hour * time.Duration(expRefreshTime)).Unix()

		access_token, err := common.GenerateJWT([]byte(jwtSecret), userID, expSecret)
		refresh_token, err := common.GenerateJWT([]byte(jwtRefreshSecret), userID, expRefresh)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_TOKEN"))
			return
		}

		h.service.CreateRefreshToken(c.Request.Context(), &models.RefreshTokenCreation{
			Token:  refresh_token,
			UserId: userID,
			Exp:    time.Unix(expRefresh, 0),
		})

		token := map[string]string{
			"access_token":  access_token,
			"refresh_token": refresh_token,
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponseWithToken(nil, token))
	}
}

func (h *handler) Logout() func(*gin.Context) {
	return func(c *gin.Context) {
		var refreshToken models.RefreshTokenSelection

		if err := c.ShouldBindJSON(&refreshToken); err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_REQUEST"))
			return
		}

		err := h.service.Logout(c.Request.Context(), refreshToken.Token)

		if err != nil {
			if customErr, ok := err.(*common.CustomError); ok {
				c.JSON(customErr.StatusCode, err)
			} else {
				c.JSON(http.StatusInternalServerError, common.NewFullCustomError(http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR"))
			}
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse("LOGOUT_SUCCESS"))

	}
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

		expSecret := time.Now().Add(time.Minute * time.Duration(expSecretTime)).Unix()
		expRefresh := time.Now().Add(time.Hour * time.Duration(expRefreshTime)).Unix()

		access_token, err := common.GenerateJWT([]byte(jwtSecret), data.ID.Hex(), expSecret)
		refresh_token, err := common.GenerateJWT([]byte(jwtRefreshSecret), data.ID.Hex(), expRefresh)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_TOKEN"))
			return
		}

		h.service.CreateRefreshToken(c.Request.Context(), &models.RefreshTokenCreation{
			Token:  refresh_token,
			UserId: data.ID.Hex(),
			Exp:    time.Unix(expRefresh, 0),
		})

		token := map[string]string{
			"access_token":  access_token,
			"refresh_token": refresh_token,
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

		expSecret := time.Now().Add(time.Minute * time.Duration(expSecretTime)).Unix()
		expRefresh := time.Now().Add(time.Hour * time.Duration(expRefreshTime)).Unix()

		access_token, err := common.GenerateJWT([]byte(jwtSecret), data.ID.Hex(), expSecret)
		refresh_token, err := common.GenerateJWT([]byte(jwtRefreshSecret), data.ID.Hex(), expRefresh)

		if err != nil {
			c.JSON(http.StatusBadRequest, common.NewFullCustomError(http.StatusBadRequest, err.Error(), "INVALID_TOKEN"))
			return
		}

		h.service.CreateRefreshToken(c.Request.Context(), &models.RefreshTokenCreation{
			Token:  refresh_token,
			UserId: data.ID.Hex(),
			Exp:    time.Unix(expRefresh, 0),
		})

		token := map[string]string{
			"access_token":  access_token,
			"refresh_token": refresh_token,
		}

		c.JSON(http.StatusCreated, common.SimpleSuccessResponseWithToken(data, token))
	}
}
