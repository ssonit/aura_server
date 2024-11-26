package utils

import (
	"context"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserModel, error)
	GetUserByID(ctx context.Context, id string) (*models.UserModel, error)
	CreateUser(ctx context.Context, user *models.UserCreation) (primitive.ObjectID, error)
	UpdateUser(ctx context.Context, id string, user *models.UserUpdate) error
	CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error
	DeleteRefreshToken(ctx context.Context, refresh_token string) error
	ListUsers(ctx context.Context, paging *common.Paging) ([]*models.UserModel, error)
}

type UserService interface {
	CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error
	Register(ctx context.Context, user *models.UserCreation) (*models.UserModel, error)
	Login(ctx context.Context, email, password string) (*models.UserModel, error)
	Logout(ctx context.Context, refresh_token string) error
	GetUser(ctx context.Context, id string) (*models.UserModel, error)
	UpdateUser(ctx context.Context, id string, user *models.UserUpdate) error
	ListUsers(ctx context.Context, paging *common.Paging) ([]*models.UserModel, error)
}
