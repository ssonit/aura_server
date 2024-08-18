package utils

import (
	"context"
	"time"

	"github.com/ssonit/aura_server/internal/auth/models"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next UserService
}

func NewLoggingMiddleware(next UserService) UserService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) Register(ctx context.Context, user *models.UserCreation) (*models.User, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Register user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.Register(ctx, user)
}

func (s *LoggingMiddleware) Login(ctx context.Context, email, password string) (*models.User, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Login user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.Login(ctx, email, password)
}
