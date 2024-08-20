package utils

import (
	"context"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next MediaService
}

func NewLoggingMiddleware(next MediaService) MediaService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) UploadImage(ctx context.Context, f *multipart.FileHeader) (primitive.ObjectID, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Upload image", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UploadImage(ctx, f)
}
