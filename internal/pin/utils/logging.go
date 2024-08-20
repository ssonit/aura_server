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
		zap.L().Info("Create pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListPinItem(ctx, filter, paging)

}
