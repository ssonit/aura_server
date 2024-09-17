package utils

import (
	"context"

	"github.com/ssonit/aura_server/internal/board/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService interface {
	CreateBoard(context.Context, *models.BoardCreation) (primitive.ObjectID, error)
	ListBoardItem(context.Context, *models.Filter) ([]models.BoardModel, error)
	GetBoardItem(context.Context, primitive.ObjectID) (*models.BoardModel, error)
	UpdateBoardItem(context.Context, string, *models.BoardUpdate) error
	SoftDeleteBoard(ctx context.Context, id string) error
	RestoreBoard(ctx context.Context, id string) error
	ListDeletedBoards(ctx context.Context, userId string) ([]models.BoardModel, error)
}

type BoardStore interface {
	CreateBoard(context.Context, *models.BoardCreation) (primitive.ObjectID, error)
	ListBoardItem(context.Context, *models.Filter) ([]models.BoardModel, error)
	GetBoardItem(context.Context, primitive.ObjectID) (*models.BoardModel, error)
	UserHasBoards(context.Context, primitive.ObjectID) (bool, error)
	UpdateBoardItem(context.Context, primitive.ObjectID, *models.BoardUpdate) error
	SoftDeleteBoard(ctx context.Context, id primitive.ObjectID) error
	RestoreBoard(ctx context.Context, id primitive.ObjectID) error
	ListDeletedBoards(ctx context.Context, userId primitive.ObjectID) ([]models.BoardModel, error)
}
