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

func (s *LoggingMiddleware) GetUser(ctx context.Context, id string) (*models.User, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Get user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetUser(ctx, id)
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

func (s *LoggingMiddleware) Logout(ctx context.Context, refresh_token string) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Logout user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.Logout(ctx, refresh_token)
}

func (s *LoggingMiddleware) CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Create refresh token", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateRefreshToken(ctx, p)
}
