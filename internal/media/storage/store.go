package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db: db}
}

func (s *store) UploadImage(ctx context.Context) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

func (s *store) GetAllImages(ctx context.Context) (interface{}, error) {
	return nil, nil
}
