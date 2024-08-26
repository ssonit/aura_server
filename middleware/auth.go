package middleware

import (
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ssonit/aura_server/common"
)

var (
	jwtSecret = common.EnvConfig("JWT_SECRET", "secret")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAuthorization := c.GetHeader("Authorization")
		if headerAuthorization == "" {
			c.JSON(http.StatusUnauthorized, common.NewFullCustomError(http.StatusUnauthorized, common.AuthorizationTokenRequired.Error(), "AUTHORIZATION_TOKEN_REQUIRED"))
			c.Abort()
			return
		}

		bearerToken := strings.Split(headerAuthorization, " ")
		tokenString := html.EscapeString(bearerToken[1])

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, common.NewFullCustomError(http.StatusUnauthorized, err.Error(), "INVALID_TOKEN"))
			c.Abort()
			return
		}

		c.Set("userID", claims["userID"])
		c.Next()
	}
}
