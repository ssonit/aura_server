package storage

import (
	"context"

	"github.com/ssonit/aura_server/internal/media/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "aura_pins"
	CollName = "medias"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db: db}
}

func (s *store) UploadImage(ctx context.Context, media *models.MediaCreation) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	data := &models.Media{
		Url:       media.Url,
		SecureUrl: media.SecureUrl,
		PublicId:  media.PublicId,
		Format:    media.Format,
		Width:     media.Width,
		Height:    media.Height,
	}

	newUser, err := collection.InsertOne(ctx, data)

	id := newUser.InsertedID.(primitive.ObjectID)

	return id, err
}
