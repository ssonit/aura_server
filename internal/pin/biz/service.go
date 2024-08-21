package biz

import (
	"context"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	store utils.PinStore
}

func NewService(store utils.PinStore) *service {
	return &service{store: store}
}

func (s *service) CreatePin(ctx context.Context, p *models.PinCreation) (primitive.ObjectID, error) {

	data, err := s.store.Create(ctx, p)
	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateEntity
	}

	return data, nil

}

func (s *service) ListPinItem(ctx context.Context, filter *models.Filter, paging *common.Paging) ([]models.PinModel, error) {
	data, err := s.store.ListItem(ctx, filter, paging)
	if err != nil {
		return nil, utils.ErrCannotGetEntity
	}

	return data, nil

}

func (s *service) GetPinById(ctx context.Context, id string) (*models.PinModel, error) {
	oID, _ := primitive.ObjectIDFromHex(id)

	data, err := s.store.GetItem(ctx, map[string]interface{}{"_id": oID})

	if err != nil {
		return nil, utils.ErrCannotGetEntity
	}

	return data, nil
}
