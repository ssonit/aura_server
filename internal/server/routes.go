package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/middleware"
	"go.uber.org/zap"

	limits "github.com/gin-contrib/size"
	media_biz "github.com/ssonit/aura_server/internal/media/biz"
	media_storage "github.com/ssonit/aura_server/internal/media/storage"
	media_http "github.com/ssonit/aura_server/internal/media/transport/gin"
	pin_biz "github.com/ssonit/aura_server/internal/pin/biz"
	pin_storage "github.com/ssonit/aura_server/internal/pin/storage"
	pin_http "github.com/ssonit/aura_server/internal/pin/transport/gin"
)

func (s *Server) MapRoutes(r *gin.Engine, httpAddr string) error {

	// Pin
	pinStore := pin_storage.NewStore(s.db)
	pinService := pin_biz.NewService(pinStore)
	pinHandler := pin_http.NewHandler(pinService)

	// Media
	mediaStore := media_storage.NewStore(s.db)
	mediaService := media_biz.NewService(mediaStore)
	mediaHandler := media_http.NewHandler(mediaService)

	// Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Accept", "X-Request-ID", "X-CSRF-Token"},
	}))
	r.Use(middleware.Recovery())
	r.Use(secure.New(secure.Config{}))
	r.Use(requestid.New())
	r.Use(limits.RequestSizeLimiter(10 * 1024 * 1024)) // 10MB

	// Routes
	pinGroup := r.Group("/pin")
	mediaGroup := r.Group("/media")

	pinHandler.RegisterRoutes(pinGroup)
	mediaHandler.RegisterRoutes(mediaGroup)

	// Health check
	r.GET("/ping", pinHandler.Ping)

	// Start server
	s.logger.Info("Server listening on ", zap.String("port", httpAddr))
	r.Run(httpAddr)

	return nil
}
