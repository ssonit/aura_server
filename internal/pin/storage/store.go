package storage

import (
	"context"
	"fmt"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"go.mongodb.org/mongo-driver/bson"
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

	if err != nil {
		return primitive.NilObjectID, utils.ErrNotInserted
	}

	id := newUser.InsertedID.(primitive.ObjectID)

	return id, nil
}

func (s *store) GetItem(ctx context.Context, filter map[string]interface{}) (*models.PinModel, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "media_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "media"},
				},
			},
		},
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
					{Key: "path", Value: "$media"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
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
		return nil, fmt.Errorf("failed to find items: %v", err)
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var item models.PinModel
		if err := cursor.Decode(&item); err != nil {
			return nil, fmt.Errorf("failed to decode item: %v", err)
		}

		// Remove the password from the user data before returning
		item.User.Password = ""

		return &item, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return nil, nil
}

func (s *store) ListItem(ctx context.Context, filter *models.Filter, paging *common.Paging, moreKeys ...string) ([]models.PinModel, error) {

	collection := s.db.Database(DbName).Collection(CollName)

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents: %v", err)
	}

	paging.Total = total

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "media_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "media"},
				},
			},
		},
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
					{Key: "path", Value: "$media"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
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
		bson.D{{Key: "$skip", Value: int64((paging.Page - 1) * paging.Limit)}},
		bson.D{{Key: "$limit", Value: int64(paging.Limit)}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, fmt.Errorf("failed to find items: %v", err)
	}
	defer cursor.Close(ctx)

	var items []models.PinModel

	if err = cursor.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("failed to decode items: %v", err)
	}

	return items, nil

}
