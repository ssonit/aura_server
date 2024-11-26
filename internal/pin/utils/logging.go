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

func (s *LoggingMiddleware) ListTags(ctx context.Context, paging *common.Paging) ([]models.Tag, error) {

	start := time.Now()

	defer func() {
		zap.L().Info("List tags", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListTags(ctx, paging)
}

func (s *LoggingMiddleware) ListSuggestions(ctx context.Context, keyword string, limit int) ([]models.Suggestion, error) {

	start := time.Now()

	defer func() {
		zap.L().Info("List suggestions", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListSuggestions(ctx, keyword, limit)
}

func (s *LoggingMiddleware) ListSoftDeletedPins(ctx context.Context, userId string) ([]models.PinModel, error) {

	start := time.Now()

	defer func() {
		zap.L().Info("List soft deleted pins", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListSoftDeletedPins(ctx, userId)
}

func (s *LoggingMiddleware) RestorePin(ctx context.Context, id, userId string) error {

	start := time.Now()

	defer func() {
		zap.L().Info("Restore pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.RestorePin(ctx, id, userId)

}

func (s *LoggingMiddleware) SoftDeletePin(ctx context.Context, id, userId string) error {

	start := time.Now()

	defer func() {
		zap.L().Info("Soft delete pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.SoftDeletePin(ctx, id, userId)
}

func (s *LoggingMiddleware) UnSaveBoardPin(ctx context.Context, p *models.BoardPinUnSave) error {

	start := time.Now()

	defer func() {
		zap.L().Info("Delete board pin item", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UnSaveBoardPin(ctx, p)
}

func (s *LoggingMiddleware) DeleteComment(ctx context.Context, commentId, userId string) error {

	start := time.Now()

	defer func() {
		zap.L().Info("Delete comment", zap.Duration("took", time.Since(start)))
	}()

	return s.next.DeleteComment(ctx, commentId, userId)
}

func (s *LoggingMiddleware) ListCommentsByPinId(ctx context.Context, pinId string, paging *common.Paging) ([]models.CommentModel, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("List comments by pin id", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListCommentsByPinId(ctx, pinId, paging)
}

func (s *LoggingMiddleware) CreateComment(ctx context.Context, comment *models.CommentCreation) (primitive.ObjectID, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("Create comment", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateComment(ctx, comment)
}

func (s *LoggingMiddleware) LikePin(ctx context.Context, like *models.LikeCreation) error {
	start := time.Now()

	defer func() {
		zap.L().Info("Like pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.LikePin(ctx, like)
}

func (s *LoggingMiddleware) UnLikePin(ctx context.Context, like *models.LikeDelete) error {
	start := time.Now()

	defer func() {
		zap.L().Info("Unlike pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UnLikePin(ctx, like)
}

func (s *LoggingMiddleware) SaveBoardPin(ctx context.Context, p *models.BoardPinSave) (primitive.ObjectID, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("Create board pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.SaveBoardPin(ctx, p)
}

func (s *LoggingMiddleware) GetBoardPinItem(ctx context.Context, filter *models.BoardPinFilter) (*models.BoardPinModel, error) {
	start := time.Now()

	defer func() {
		zap.L().Info("Get board pin item", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetBoardPinItem(ctx, filter)
}

func (s *LoggingMiddleware) ListBoardPinItem(ctx context.Context, filter *models.BoardPinFilter, paging *common.Paging) ([]models.BoardPinModel, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("List board pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListBoardPinItem(ctx, filter, paging)
}

func (s *LoggingMiddleware) UpdatePin(ctx context.Context, id string, pin *models.PinUpdate, userId string) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Update pin", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdatePin(ctx, id, pin, userId)
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

func (s *LoggingMiddleware) GetPinById(ctx context.Context, id, userId string) (*models.PinModel, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("Get pin by id", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetPinById(ctx, id, userId)

}
