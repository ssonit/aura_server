package storage

import (
	"context"

	"github.com/ssonit/aura_server/internal/board/models"
	"github.com/ssonit/aura_server/internal/board/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "aura_pins"
	CollName = "boards"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) CreateBoard(ctx context.Context, p *models.BoardCreation) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	data := &models.Board{
		UserId:    p.UserId,
		Name:      p.Name,
		IsPrivate: p.IsPrivate,
	}

	newData, err := collection.InsertOne(ctx, data)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateEntity
	}

	id := newData.InsertedID.(primitive.ObjectID)

	return id, nil
}

func (s *store) ListBoardItem(ctx context.Context, filter *models.Filter) ([]models.BoardModel, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "users"},
					{Key: "localField", Value: "user_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "user"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$user"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{
						Key:   "user.password",
						Value: 0,
					},
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, utils.ErrFailedToFindEntity
	}
	defer cursor.Close(ctx)

	var items []models.BoardModel

	if err = cursor.All(ctx, &items); err != nil {
		return nil, utils.ErrFailedToDecode
	}

	return items, nil
}
