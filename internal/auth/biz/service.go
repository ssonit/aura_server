package biz

import (
	"context"

	"github.com/ssonit/aura_server/common"
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

func (s *service) Register(ctx context.Context, user *models.UserCreation) (*models.User, error) {
	// check user exists
	_, err := s.store.GetUserByEmail(ctx, user.Email)

	if err == nil {
		return nil, common.NewCustomError(utils.ErrEmailAlreadyExists, utils.ErrEmailAlreadyExists.Error(), "EMAIL_ALREADY_EXISTS")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, utils.ErrFailedToHashPassword
	}

	user.Password = string(hashedPassword)

	id, err := s.store.CreateUser(ctx, user)

	if err != nil {
		return nil, common.ErrCannotCreateEntity("user", err)
	}

	data, err := s.store.GetUserByID(ctx, id.Hex())

	if err != nil {
		return nil, common.ErrCannotGetEntity("user", err)
	}

	data.Password = ""

	return data, nil
}

func (s *service) Login(ctx context.Context, email, password string) (*models.User, error) {
	data, err := s.store.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, common.ErrCannotGetEntity("user", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))

	if err != nil {
		return nil, common.NewCustomError(utils.ErrInvalidEmailOrPassword, utils.ErrInvalidEmailOrPassword.Error(), "INVALID_PASSWORD")
	}

	data.Password = ""

	return data, nil

}
