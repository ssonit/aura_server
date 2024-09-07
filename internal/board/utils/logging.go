package utils

import (
	"context"
	"time"

	"github.com/ssonit/aura_server/internal/board/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next BoardService
}

func NewLoggingMiddleware(next BoardService) BoardService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) GetBoardItem(ctx context.Context, id primitive.ObjectID) (*models.BoardModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("Get board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetBoardItem(ctx, id)
}

func (s *LoggingMiddleware) CreateBoard(ctx context.Context, p *models.BoardCreation) (primitive.ObjectID, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("Create board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateBoard(ctx, p)
}

func (s *LoggingMiddleware) ListBoardItem(ctx context.Context, filter *models.Filter) ([]models.BoardModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("List board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListBoardItem(ctx, filter)
}
