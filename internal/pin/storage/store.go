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
	DbName           = "aura_pins"
	CollName         = "pins"
	CollNameBoardPin = "board_pins"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) DeleteBoardPin(ctx context.Context, filter *models.BoardPinFilter) error {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	fmt.Println(filter.BoardId, filter.PinId, "filter")
	_, err := collection.DeleteOne(ctx, bson.M{"pin_id": filter.PinId, "user_id": filter.UserId})

	if err != nil {
		return utils.ErrCannotDeleteBoardPin
	}

	return nil
}

func (s *store) GetBoardPinItem(ctx context.Context, filter *models.BoardPinFilter) (*models.BoardPinModel, error) {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "pin_id", Value: filter.PinId}, {Key: "user_id", Value: filter.UserId}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "pins"},
					{Key: "localField", Value: "pin_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "pin"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "boards"},
					{Key: "localField", Value: "board_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "board"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "pin.media_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "media"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$pin"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$board"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
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
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {

		return nil, utils.ErrFailedToFindEntity
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var item models.BoardPinModel
		if err := cursor.Decode(&item); err != nil {
			return nil, utils.ErrFailedToDecode
		}

		return &item, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, utils.ErrCursorError
	}

	return nil, nil
}

func (s *store) ListBoardPinItem(ctx context.Context, filter *models.BoardPinFilter, paging *common.Paging) ([]models.BoardPinModel, error) {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	total, err := collection.CountDocuments(ctx, bson.D{{Key: "board_id", Value: filter.BoardId}})
	if err != nil {
		return nil, utils.ErrFailedToCount
	}

	paging.Total = total

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "board_id", Value: filter.BoardId}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "pins"},
					{Key: "localField", Value: "pin_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "pin"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "boards"},
					{Key: "localField", Value: "board_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "board"},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "pin.media_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "media"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$pin"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$board"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
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
		bson.D{{Key: "$skip", Value: int64((paging.Page - 1) * paging.Limit)}},
		bson.D{{Key: "$limit", Value: int64(paging.Limit)}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, utils.ErrFailedToFindEntity
	}
	defer cursor.Close(ctx)

	var items []models.BoardPinModel

	if err = cursor.All(ctx, &items); err != nil {
		return nil, utils.ErrFailedToDecode
	}

	return items, nil

}

func (s *store) CreateBoardPin(ctx context.Context, p *models.BoardPinCreation) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	data := &models.BoardPin{
		BoardId: p.BoardId,
		PinId:   p.PinId,
		UserId:  p.UserId,
	}

	newData, err := collection.InsertOne(ctx, data)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateBoardPin
	}

	id := newData.InsertedID.(primitive.ObjectID)

	return id, nil
}

func (s *store) UpdatePin(ctx context.Context, id string, pin *models.PinUpdate) error {
	collection := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(id)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: pin.Title}, {Key: "description", Value: pin.Description}, {Key: "link_url", Value: pin.LinkUrl}}}}

	_, err := collection.UpdateByID(ctx, oID, update)

	if err != nil {
		return utils.ErrCannotUpdatePin
	}

	return nil
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

	newData, err := collection.InsertOne(ctx, data)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreatePin
	}

	id := newData.InsertedID.(primitive.ObjectID)

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
		return nil, utils.ErrFailedToFindEntity
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var item models.PinModel
		if err := cursor.Decode(&item); err != nil {
			return nil, utils.ErrFailedToDecode
		}

		// Remove the password from the user data before returning
		item.User.Password = ""

		return &item, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, utils.ErrCursorError
	}

	return nil, nil
}

func (s *store) ListItem(ctx context.Context, filter *models.Filter, paging *common.Paging, moreKeys ...string) ([]models.PinModel, error) {

	collection := s.db.Database(DbName).Collection(CollName)

	user_id, _ := primitive.ObjectIDFromHex(filter.UserId)

	filterMap := map[string]interface{}{}

	if filter.UserId != "" {
		filterMap["user_id"] = user_id
	}

	if filter.Title != "" {
		filterMap["title"] = filter.Title
	}

	total, err := collection.CountDocuments(ctx, filterMap)
	if err != nil {
		return nil, utils.ErrFailedToCount
	}

	paging.Total = total

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filterMap}},
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
		return nil, utils.ErrFailedToFindEntity
	}
	defer cursor.Close(ctx)

	var items []models.PinModel

	if err = cursor.All(ctx, &items); err != nil {
		return nil, utils.ErrFailedToDecode
	}

	return items, nil

}
