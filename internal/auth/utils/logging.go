package utils

import (
	"context"
	"time"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next UserService
}

func NewLoggingMiddleware(next UserService) UserService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) UnbannedUser(ctx context.Context, id string) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Unbanned user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UnbannedUser(ctx, id)
}

func (s *LoggingMiddleware) BannedUser(ctx context.Context, id string) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Banned user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.BannedUser(ctx, id)
}

func (s *LoggingMiddleware) ListUsers(ctx context.Context, paging *common.Paging) ([]*models.UserModel, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("List users", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ListUsers(ctx, paging)
}

func (s *LoggingMiddleware) UpdateUser(ctx context.Context, id string, user *models.UserUpdate) error {
	start := time.Now()
	defer func() {
		zap.L().Info("Update user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdateUser(ctx, id, user)
}

func (s *LoggingMiddleware) GetUser(ctx context.Context, id string) (*models.UserModel, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Get user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetUser(ctx, id)
}

func (s *LoggingMiddleware) Register(ctx context.Context, user *models.UserCreation) (*models.UserModel, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("Register user", zap.Duration("took", time.Since(start)))
	}()

	return s.next.Register(ctx, user)
}

func (s *LoggingMiddleware) Login(ctx context.Context, email, password string) (*models.UserModel, error) {
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
