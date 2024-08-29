package biz

import (
	"context"

	"github.com/ssonit/aura_server/internal/auth/models"
	"github.com/ssonit/aura_server/internal/auth/utils"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	store utils.UserStore
}

func NewService(s utils.UserStore) *service {
	return &service{store: s}
}

func (s *service) CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error {
	return s.store.CreateRefreshToken(ctx, p)
}

func (s *service) Logout(ctx context.Context, refresh_token string) error {
	return s.store.DeleteRefreshToken(ctx, refresh_token)
}

func (s *service) Register(ctx context.Context, user *models.UserCreation) (*models.User, error) {
	// check user exists
	_, err := s.store.GetUserByEmail(ctx, user.Email)

	if err == nil {
		return nil, utils.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, utils.ErrFailedToHashPassword
	}

	user.Password = string(hashedPassword)

	id, err := s.store.CreateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	data, err := s.store.GetUserByID(ctx, id.Hex())

	if err != nil {
		return nil, err
	}

	data.Password = ""

	return data, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*models.User, error) {
	data, err := s.store.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))

	if err != nil {
		return nil, utils.ErrEmailOrPassInvalid
	}

	data.Password = ""

	return data, nil

}
