package utils

import (
	"context"

	"github.com/ssonit/aura_server/internal/board/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService interface {
	CreateBoard(context.Context, *models.BoardCreation) (primitive.ObjectID, error)
	ListBoardItem(context.Context, *models.Filter) ([]models.BoardModel, error)
}

type BoardStore interface {
	CreateBoard(context.Context, *models.BoardCreation) (primitive.ObjectID, error)
	ListBoardItem(context.Context, *models.Filter) ([]models.BoardModel, error)
}
