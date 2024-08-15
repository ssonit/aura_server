package utils

import (
	"context"
	"time"

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
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreatePin(ctx, p)

}
