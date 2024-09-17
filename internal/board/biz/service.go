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

func (s *service) ListDeletedBoards(ctx context.Context, userId string) ([]models.BoardModel, error) {
	OId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, utils.ErrFailedToDecodeObjID
	}

	return s.store.ListDeletedBoards(ctx, OId)
}

func (s *service) RestoreBoard(ctx context.Context, id string) error {
	boardId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	return s.store.RestoreBoard(ctx, boardId)
}

func (s *service) SoftDeleteBoard(ctx context.Context, id string) error {
	boardId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	return s.store.SoftDeleteBoard(ctx, boardId)
}

func (s *service) UpdateBoardItem(ctx context.Context, id string, p *models.BoardUpdate) error {
	boardId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrFailedToDecodeObjID
	}

	return s.store.UpdateBoardItem(ctx, boardId, p)

}

func (s *service) GetBoardItem(ctx context.Context, id primitive.ObjectID) (*models.BoardModel, error) {
	data, err := s.store.GetBoardItem(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *service) CreateBoard(ctx context.Context, p *models.BoardCreation) (primitive.ObjectID, error) {
	userHasBoards, err := s.store.UserHasBoards(ctx, p.UserId)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if !userHasBoards {
		_, err := s.store.CreateBoard(ctx, &models.BoardCreation{
			UserId:    p.UserId,
			Name:      "All Pins",
			IsPrivate: true,
			Type:      "all_pins",
		})

		if err != nil {
			return primitive.NilObjectID, err
		}
	}

	data, err := s.store.CreateBoard(ctx, &models.BoardCreation{
		UserId:    p.UserId,
		Name:      p.Name,
		IsPrivate: p.IsPrivate,
		Type:      "custom",
	})
	if err != nil {
		return primitive.NilObjectID, err
	}

	return data, nil
}

func (s *service) ListBoardItem(ctx context.Context, filter *models.Filter) ([]models.BoardModel, error) {
	data, err := s.store.ListBoardItem(ctx, filter)
	if err != nil {
		return nil, err
	}

	return data, nil
}
