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
	GetPinById(ctx context.Context, id, userId string) (*models.PinModel, error)
	ListBoardPinItem(context.Context, *models.BoardPinFilter, *common.Paging) ([]models.BoardPinModel, error)
	GetBoardPinItem(context.Context, *models.BoardPinFilter) (*models.BoardPinModel, error)
	UpdatePin(ctx context.Context, id string, pin *models.PinUpdate, userId string) error
	SaveBoardPin(context.Context, *models.BoardPinSave) (primitive.ObjectID, error)
	LikePin(context.Context, *models.LikeCreation) error
	UnLikePin(context.Context, *models.LikeDelete) error
	CreateComment(context.Context, *models.CommentCreation) (primitive.ObjectID, error)
	ListCommentsByPinId(context.Context, string, *common.Paging) ([]models.CommentModel, error)
	DeleteComment(context.Context, string, string) error
	UnSaveBoardPin(context.Context, *models.BoardPinUnSave) error
	SoftDeletePin(ctx context.Context, id, userId string) error
	RestorePin(ctx context.Context, id, userId string) error
	ListSoftDeletedPins(context.Context, string) ([]models.PinModel, error)
	ListSuggestions(ctx context.Context, keyword string, limit int) ([]models.Suggestion, error)
	ListTags(context.Context, *common.Paging) ([]models.Tag, error)
	CreateTag(context.Context, models.TagCreation) (primitive.ObjectID, error)
	DeleteTag(context.Context, string) error
}

type PinStore interface {
	Create(context.Context, *models.PinCreation, []primitive.ObjectID) (primitive.ObjectID, error)
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
	LikePin(ctx context.Context, userID, pinID primitive.ObjectID) error
	UnlikePin(ctx context.Context, userID, pinID primitive.ObjectID) error
	CreateComment(context.Context, *models.CommentCreationStore) (primitive.ObjectID, error)
	ListCommentsByPinId(context.Context, primitive.ObjectID, *common.Paging) ([]models.CommentModel, error)
	DeleteComment(context.Context, primitive.ObjectID) error
	GetCommentById(context.Context, primitive.ObjectID) (*models.CommentModel, error)
	SoftDeletePin(context.Context, primitive.ObjectID) error
	RestorePin(context.Context, primitive.ObjectID) error
	ListSoftDeletedPins(context.Context, primitive.ObjectID) ([]models.PinModel, error)
	CheckAndCreateTags(context.Context, []string) ([]primitive.ObjectID, error)
	MatchingTags(ctx context.Context, keyword string) ([]primitive.ObjectID, error)
	CheckAndCreateSuggestions(context.Context, []string) ([]primitive.ObjectID, error)
	ListSuggestions(ctx context.Context, keyword string, limit int) ([]models.Suggestion, error)
	ListTags(context.Context, *common.Paging) ([]models.Tag, error)
	DeleteTag(context.Context, primitive.ObjectID) error
}
