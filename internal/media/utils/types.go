package utils

import (
	"context"
	"mime/multipart"

	"github.com/ssonit/aura_server/internal/media/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService interface {
	UploadImage(context.Context, *multipart.FileHeader) (primitive.ObjectID, error)
	GetMedia(ctx context.Context, id string) (*models.Media, error)
}

type MediaStore interface {
	UploadImage(context.Context, *models.MediaCreation) (primitive.ObjectID, error)
	GetMedia(ctx context.Context, id string) (*models.Media, error)
}
