package server

import (
	"github.com/gin-gonic/gin"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/utils"

	pin_http "github.com/ssonit/aura_server/internal/pin/transport/gin"
)

var (
	httpAddr = common.EnvConfig("HTTP_ADDR", ":8080")
)

type Server struct {
	pin_handler utils.PinHandler
}

func NewServer() *Server {
	handler := pin_http.NewHandler()
	return &Server{
		pin_handler: handler,
	}
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", s.pin_handler.Ping)
}
