package utils

import (
	"context"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PinService interface {
	CreatePin(context.Context, *models.PinCreation) (primitive.ObjectID, error)
	ListPinItem(ctx context.Context, filter *models.Filter, paging *common.Paging) ([]models.PinModel, error)
	GetPinById(ctx context.Context, id string) (*models.PinModel, error)
}

type PinStore interface {
	Create(context.Context, *models.PinCreation) (primitive.ObjectID, error)
	ListItem(ctx context.Context, filter *models.Filter, paging *common.Paging, moreKeys ...string) ([]models.PinModel, error)
	GetItem(ctx context.Context, filter map[string]interface{}) (*models.PinModel, error)
}
