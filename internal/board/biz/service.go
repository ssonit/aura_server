package biz

import (
	"context"

	"github.com/ssonit/aura_server/internal/board/models"
	"github.com/ssonit/aura_server/internal/board/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	store utils.BoardStore
}

func NewService(store utils.BoardStore) *service {
	return &service{store: store}
}

func (s *service) CreateBoard(ctx context.Context, p *models.BoardCreation) (primitive.ObjectID, error) {
	data, err := s.store.CreateBoard(ctx, p)
	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateEntity
	}

	return data, nil
}

func (s *service) ListBoardItem(ctx context.Context, filter *models.Filter) ([]models.BoardModel, error) {
	data, err := s.store.ListBoardItem(ctx, filter)
	if err != nil {
		return nil, utils.ErrCannotGetEntity
	}

	return data, nil
}
