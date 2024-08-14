package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Server struct {
	r      *gin.Engine
	db     *mongo.Client
	logger *zap.Logger
}

func NewServer(r *gin.Engine, db *mongo.Client, logger *zap.Logger) *Server {
	return &Server{
		r:      r,
		db:     db,
		logger: logger,
	}
}

func (s *Server) Run(httpAddr string) error {
	if err := s.MapRoutes(s.r, httpAddr); err != nil {
		return err
	}

	return nil
}
