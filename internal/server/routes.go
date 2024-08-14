package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/internal/pin/biz"
	"github.com/ssonit/aura_server/internal/pin/storage"
	"github.com/ssonit/aura_server/middleware"
	"go.uber.org/zap"

	limits "github.com/gin-contrib/size"
	pin_http "github.com/ssonit/aura_server/internal/pin/transport/gin"
)

func (s *Server) MapRoutes(r *gin.Engine, httpAddr string) error {
	pinStore := storage.NewStore(s.db)
	pinService := biz.NewService(pinStore)
	pinHandler := pin_http.NewHandler(pinService)

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Accept", "X-Request-ID", "X-CSRF-Token"},
	}))
	r.Use(middleware.Recovery())
	r.Use(secure.New(secure.Config{}))
	r.Use(requestid.New())
	r.Use(limits.RequestSizeLimiter(10))

	r.GET("/ping", pinHandler.Ping)

	s.logger.Info("Server listening on ", zap.String("port", httpAddr))
	r.Run(httpAddr)

	return nil
}
