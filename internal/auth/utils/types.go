package utils

import (
	"context"

	"github.com/ssonit/aura_server/internal/auth/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.UserCreation) (primitive.ObjectID, error)
	CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error
	DeleteRefreshToken(ctx context.Context, refresh_token string) error
}

type UserService interface {
	CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error
	Register(ctx context.Context, user *models.UserCreation) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
	Logout(ctx context.Context, refresh_token string) error
}
