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

func (s *LoggingMiddleware) ListDeletedBoards(ctx context.Context, userId string) ([]models.BoardModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("List deleted boards", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListDeletedBoards(ctx, userId)
}

func (s *LoggingMiddleware) RestoreBoard(ctx context.Context, id string) error {

	start := time.Now()

	defer func() {
		zap.L().Info("Restore board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.RestoreBoard(ctx, id)
}

func (s *LoggingMiddleware) SoftDeleteBoard(ctx context.Context, id string) error {

	start := time.Now()
	defer func() {
		zap.L().Info("Soft delete board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.SoftDeleteBoard(ctx, id)
}

func (s *LoggingMiddleware) UpdateBoardItem(ctx context.Context, id string, p *models.BoardUpdate) error {

	start := time.Now()
	defer func() {
		zap.L().Info("Update board", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdateBoardItem(ctx, id, p)
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
