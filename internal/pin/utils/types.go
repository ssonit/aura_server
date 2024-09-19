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
	ListBoardPinItem(context.Context, *models.BoardPinFilter, *common.Paging) ([]models.BoardPinModel, error)
	GetBoardPinItem(context.Context, *models.BoardPinFilter) (*models.BoardPinModel, error)
	UpdatePin(ctx context.Context, id string, pin *models.PinUpdate, userId string) error
	SaveBoardPin(context.Context, *models.BoardPinSave) (primitive.ObjectID, error)
}

type PinStore interface {
	Create(context.Context, *models.PinCreation) (primitive.ObjectID, error)
	ListItem(ctx context.Context, filter *models.Filter, paging *common.Paging, moreKeys ...string) ([]models.PinModel, error)
	GetItem(ctx context.Context, filter map[string]interface{}) (*models.PinModel, error)
	UpdatePin(ctx context.Context, id string, pin *models.PinUpdate) error
	CreateBoardPin(context.Context, *models.BoardPinCreation) (primitive.ObjectID, error)
	ListBoardPinItem(context.Context, *models.BoardPinFilter, *common.Paging) ([]models.BoardPinModel, error)
	GetBoardPinItem(context.Context, *models.BoardPinFilter) (*models.BoardPinModel, error)
	DeleteBoardPin(context.Context, *models.BoardPinFilter) error
	GetBoardByUserId(context.Context, primitive.ObjectID, string) (primitive.ObjectID, error)
	DeleteBoardPinById(context.Context, primitive.ObjectID) error
	IsPinOwnedByUser(ctx context.Context, userId, pinId primitive.ObjectID) (bool, error)
	CheckIfPinExistsInBoard(ctx context.Context, boardId primitive.ObjectID, pinId primitive.ObjectID) (bool, error)
}
