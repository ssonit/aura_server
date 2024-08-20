package storage

import (
	"context"

	"github.com/ssonit/aura_server/internal/pin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "aura_pins"
	CollName = "pins"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, p *models.PinCreation) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	data := &models.Pin{
		Title:       p.Title,
		Description: p.Description,
		UserId:      p.UserId,
		MediaId:     p.MediaId,
		LinkUrl:     p.LinkUrl,
	}

	newUser, err := collection.InsertOne(ctx, data)

	id := newUser.InsertedID.(primitive.ObjectID)

	return id, err
}
