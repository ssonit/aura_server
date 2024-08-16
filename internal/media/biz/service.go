package biz

import (
	"context"

	"github.com/ssonit/aura_server/internal/media/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	store utils.MediaStore
}

func NewService(store utils.MediaStore) *service {
	return &service{
		store: store,
	}
}

// UploadImage uploads an image to the server
func (s *service) UploadImage(ctx context.Context) (primitive.ObjectID, error) {
	return s.store.UploadImage(ctx)
}

// GetAllImages returns all images from the server
func (s *service) GetAllImages(ctx context.Context) (interface{}, error) {
	return s.store.GetAllImages(ctx)
}
