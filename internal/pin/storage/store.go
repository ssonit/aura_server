package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/pin/models"
	"github.com/ssonit/aura_server/internal/pin/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	Board_Model "github.com/ssonit/aura_server/internal/board/models"
)

const (
	DbName           = "aura_pins"
	CollName         = "pins"
	CollNameBoardPin = "board_pins"
	CollNameBoard    = "boards"
	CollNameLikes    = "likes"
	CollNameComments = "comments"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) ListSoftDeletedPins(ctx context.Context, userId primitive.ObjectID) ([]models.PinModel, error) {

	collection := s.db.Database(DbName).Collection(CollName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}, {Key: "deleted_at", Value: bson.D{{Key: "$ne", Value: nil}}}}}},
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
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$media"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "deleted_at", Value: -1}}}},
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

func (s *store) RestorePin(ctx context.Context, id primitive.ObjectID) error {

	collection := s.db.Database(DbName).Collection(CollName)

	update := bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: "deleted_at", Value: ""},
		}},
	}

	_, err := collection.UpdateByID(ctx, id, update)
	if err != nil {
		return utils.ErrFailedRestore
	}

	return nil
}

func (s *store) SoftDeletePin(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.Database(DbName).Collection(CollName)

	now := time.Now()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "deleted_at", Value: now},
		}},
	}

	_, err := collection.UpdateByID(ctx, id, update)

	if err != nil {
		return utils.ErrCannotSoftDelete
	}

	return nil

}

func (s *store) GetCommentById(ctx context.Context, id primitive.ObjectID) (*models.CommentModel, error) {
	collection := s.db.Database(DbName).Collection(CollNameComments)

	var data *models.CommentModel

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)

	if err != nil {
		return nil, utils.ErrCommentNotExists
	}

	return data, nil
}

func (s *store) DeleteComment(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.Database(DbName).Collection(CollNameComments)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return utils.ErrCannotDeleteComment
	}

	return nil
}

func (s *store) ListCommentsByPinId(ctx context.Context, pinId primitive.ObjectID, paging *common.Paging) ([]models.CommentModel, error) {
	collection := s.db.Database(DbName).Collection(CollNameComments)

	total, err := collection.CountDocuments(ctx, bson.D{{Key: "pin_id", Value: pinId}})

	if err != nil {
		return nil, utils.ErrFailedToCount
	}

	paging.Total = total

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "pin_id", Value: pinId}}}},
		bson.D{{Key: "$skip", Value: int64((paging.Page - 1) * paging.Limit)}},
		bson.D{{Key: "$limit", Value: int64(paging.Limit)}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}, // Sắp xếp theo ngày tạo mới nhất
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
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "user.avatar_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "user.avatar"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$user.avatar"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "user.password", Value: 0}, // Loại bỏ trường password của user
				},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, utils.ErrFailedToFindEntity
	}

	defer cursor.Close(ctx)

	var items []models.CommentModel

	if err = cursor.All(ctx, &items); err != nil {
		return nil, utils.ErrFailedToDecode

	}

	return items, nil

}

func (s *store) CreateComment(ctx context.Context, p *models.CommentCreationStore) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollNameComments)

	data := &models.Comment{
		PinId:   p.PinId,
		UserId:  p.UserId,
		Content: p.Content,
	}

	newData, err := collection.InsertOne(ctx, data)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateComment
	}

	id := newData.InsertedID.(primitive.ObjectID)

	return id, nil

}

func (s *store) UnlikePin(ctx context.Context, userID, pinID primitive.ObjectID) error {
	collection := s.db.Database(DbName).Collection(CollNameLikes)

	_, err := collection.DeleteOne(ctx, bson.M{"user_id": userID, "pin_id": pinID})

	if err != nil {
		return utils.ErrCannotUnlikePin
	}

	return nil
}

func (s *store) LikePin(ctx context.Context, userID, pinID primitive.ObjectID) error {
	collection := s.db.Database(DbName).Collection(CollNameLikes)

	count, err := collection.CountDocuments(ctx, bson.M{"user_id": userID, "pin_id": pinID})
	if err != nil {
		return err
	}

	if count > 0 {
		return utils.ErrAlreadyLiked
	}

	like := &models.Like{
		UserId: userID,
		PinId:  pinID,
	}

	_, err = collection.InsertOne(ctx, like)

	if err != nil {
		return utils.ErrCannotLikePin
	}

	return nil

}

func (s *store) CheckIfPinExistsInBoard(ctx context.Context, boardId primitive.ObjectID, pinId primitive.ObjectID) (bool, error) {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	filter := bson.M{
		"board_id": boardId,
		"pin_id":   pinId,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *store) IsPinOwnedByUser(ctx context.Context, userId, pinId primitive.ObjectID) (bool, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	filter := bson.M{
		"user_id": userId,
		"_id":     pinId,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *store) GetBoardByUserId(ctx context.Context, id primitive.ObjectID, board_type string) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollNameBoard)

	var data Board_Model.BoardModel

	if err := collection.FindOne(ctx, bson.M{"user_id": id, "type": board_type}).Decode(&data); err != nil {
		return primitive.NilObjectID, utils.ErrFailedToFindEntity
	}

	return data.ID, nil

}

func (s *store) DeleteBoardPin(ctx context.Context, filter *models.BoardPinFilter) error {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	filterMap := bson.M{
		"pin_id":  filter.PinId,
		"user_id": filter.UserId,
	}

	if !filter.BoardId.IsZero() {
		filterMap["board_id"] = filter.BoardId
	}

	_, err := collection.DeleteOne(ctx, filterMap)

	if err != nil {
		return utils.ErrCannotDeleteBoardPin
	}

	fmt.Println(filterMap)

	return nil
}

func (s *store) DeleteBoardPinById(ctx context.Context, id primitive.ObjectID) error {
	collection := s.db.Database(DbName).Collection(CollNameBoardPin)

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

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
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "board.type", Value: "custom"},
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
			{Key: "$match",
				Value: bson.D{
					{Key: "pin.deleted_at", Value: bson.M{"$eq": nil}}, // Lọc các pin chưa bị xóa
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
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
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
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: filter["_id"]}}}},
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
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "user.avatar_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "user.avatar"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$user.avatar"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "likes"},
					{Key: "let", Value: bson.D{
						{Key: "pin_id", Value: "$_id"},
					}},
					{Key: "pipeline", Value: mongo.Pipeline{
						// Match likes for the current user and pin
						bson.D{
							{Key: "$match", Value: bson.D{
								{Key: "$expr", Value: bson.D{
									{Key: "$and", Value: bson.A{
										bson.D{{Key: "$eq", Value: bson.A{"$user_id", filter["user_id"]}}},
										bson.D{{Key: "$eq", Value: bson.A{"$pin_id", "$$pin_id"}}},
									}},
								}},
							}},
						},
					}},
					{Key: "as", Value: "likes"},
				},
			},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "isLiked", Value: bson.D{
					{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$likes"}}, 0}}},
						true,
						false,
					}},
				}},
			}},
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

	filterMap := bson.M{
		"deleted_at": nil,
	}

	// filterMap := map[string]interface{}{
	// 	"deleted_at": nil,
	// }

	if filter.UserId != "" {

		filterMap["user_id"] = user_id
	}

	if filter.Title != "" {
		filterMap["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}

	var sortOrder int = 1
	if filter.Sort == "desc" {
		sortOrder = -1
	}

	var sortKey string = "created_at"
	if filter.SortKey != "" {
		sortKey = filter.SortKey
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
		bson.D{{Key: "$sort", Value: bson.D{{Key: sortKey, Value: sortOrder}}}},
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
