package utils

import "github.com/gin-gonic/gin"

type PinService interface {
}

type PinStore interface{}

type PinHandler interface {
	Ping(c *gin.Context)
}
