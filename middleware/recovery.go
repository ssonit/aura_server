package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
)

func Recovery() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrInternalServer(err))
				}

				panic(r)
			}
		}()

		c.Next()
	}
}
