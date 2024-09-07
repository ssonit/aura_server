package utils

import (
	"context"
	"time"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next PinService
}

func NewLoggingMiddleware(next PinService) PinService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) ListBoardPinItem(ctx context.Context, filter *models.BoardPinFilter, paging *common.Paging) ([]models.BoardPinModel, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("List board pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListBoardPinItem(ctx, filter, paging)
}

func (s *LoggingMiddleware) CreatePin(ctx context.Context, p *models.PinCreation) (primitive.ObjectID, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("Create pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreatePin(ctx, p)

}

func (s *LoggingMiddleware) ListPinItem(ctx context.Context, filter *models.Filter, paging *common.Paging) ([]models.PinModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("List pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListPinItem(ctx, filter, paging)

}

func (s *LoggingMiddleware) GetPinById(ctx context.Context, id string) (*models.PinModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("Get pin by id", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetPinById(ctx, id)

}
