package utils

import (
	"context"

	"github.com/ssonit/aura_server/internal/pin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PinService interface {
	CreatePin(context.Context, *models.PinCreation) (primitive.ObjectID, error)
}

type PinStore interface {
	Create(context.Context, *models.PinCreation) (primitive.ObjectID, error)
}
