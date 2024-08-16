package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService interface {
	UploadImage(context.Context) (primitive.ObjectID, error)
	GetAllImages(context.Context) (interface{}, error)
}

type MediaStore interface {
	UploadImage(context.Context) (primitive.ObjectID, error)
	GetAllImages(context.Context) (interface{}, error)
}
