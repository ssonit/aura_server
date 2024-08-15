package biz

import (
	"context"

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
	return s.store.Create(ctx, p)
}
