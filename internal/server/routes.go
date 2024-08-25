package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/middleware"
	"go.uber.org/zap"

	limits "github.com/gin-contrib/size"

	user_biz "github.com/ssonit/aura_server/internal/auth/biz"
	user_storage "github.com/ssonit/aura_server/internal/auth/storage"
	user_http "github.com/ssonit/aura_server/internal/auth/transport/gin"
	user_logging "github.com/ssonit/aura_server/internal/auth/utils"

	media_biz "github.com/ssonit/aura_server/internal/media/biz"
	media_storage "github.com/ssonit/aura_server/internal/media/storage"
	media_http "github.com/ssonit/aura_server/internal/media/transport/gin"
	media_logging "github.com/ssonit/aura_server/internal/media/utils"

	pin_biz "github.com/ssonit/aura_server/internal/pin/biz"
	pin_storage "github.com/ssonit/aura_server/internal/pin/storage"
	pin_http "github.com/ssonit/aura_server/internal/pin/transport/gin"
	pin_logging "github.com/ssonit/aura_server/internal/pin/utils"

	board_biz "github.com/ssonit/aura_server/internal/board/biz"
	board_storage "github.com/ssonit/aura_server/internal/board/storage"
	board_http "github.com/ssonit/aura_server/internal/board/transport/gin"
	board_logging "github.com/ssonit/aura_server/internal/board/utils"
)

func (s *Server) MapRoutes(r *gin.Engine, httpAddr string) error {

	// Pin
	pinStore := pin_storage.NewStore(s.db)
	pinService := pin_biz.NewService(pinStore)
	pinServiceWithLogging := pin_logging.NewLoggingMiddleware(pinService)
	pinHandler := pin_http.NewHandler(pinServiceWithLogging)

	// Media
	mediaStore := media_storage.NewStore(s.db)
	mediaService := media_biz.NewService(mediaStore)
	mediaServiceWithLogging := media_logging.NewLoggingMiddleware(mediaService)
	mediaHandler := media_http.NewHandler(mediaServiceWithLogging)

	// User
	userStore := user_storage.NewStore(s.db)
	userService := user_biz.NewService(userStore)
	userServiceWithLogging := user_logging.NewLoggingMiddleware(userService)
	userHandler := user_http.NewHandler(userServiceWithLogging)

	// Board
	boardStore := board_storage.NewStore(s.db)
	boardService := board_biz.NewService(boardStore)
	boardServiceWithLogging := board_logging.NewLoggingMiddleware(boardService)
	boardHandler := board_http.NewHandler(boardServiceWithLogging)

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
	userGroup := r.Group("/user")
	boardGroup := r.Group("/board")

	pinHandler.RegisterRoutes(pinGroup)
	mediaHandler.RegisterRoutes(mediaGroup)
	userHandler.RegisterRoutes(userGroup)
	boardHandler.RegisterRoutes(boardGroup)

	// Health check
	r.GET("/ping", middleware.AuthMiddleware(), pinHandler.Ping)

	// Start server
	s.logger.Info("Server listening on ", zap.String("port", httpAddr))
	r.Run(httpAddr)

	return nil
}
