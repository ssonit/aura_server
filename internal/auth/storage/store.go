package storage

import (
	"context"
	"fmt"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/auth/models"
	"github.com/ssonit/aura_server/internal/auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName               = "aura_pins"
	CollName             = "users"
	CollNameRefreshToken = "refresh_tokens"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db: db}
}

func (s *store) ListUsers(ctx context.Context, paging *common.Paging) ([]*models.UserModel, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	total, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, utils.ErrFailedToFindEntity
	}

	paging.Total = total

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "avatar_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "avatar"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$avatar"},
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

	var users []*models.UserModel

	for cursor.Next(ctx) {
		var item models.UserModel
		if err := cursor.Decode(&item); err != nil {
			return nil, utils.ErrFailedToDecode
		}

		users = append(users, &item)
	}

	if err := cursor.Err(); err != nil {
		return nil, utils.ErrCursorError
	}

	return users, nil
}

func (s *store) UpdateUser(ctx context.Context, id string, user *models.UserUpdate) error {

	collection := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(id)

	update := bson.D{
		{Key: "username", Value: user.Username},
		{Key: "bio", Value: user.Bio},
		{Key: "website", Value: user.Website},
	}

	if user.AvatarID != "" {
		avatarOID, _ := primitive.ObjectIDFromHex(user.AvatarID)
		update = append(update, bson.E{Key: "avatar_id", Value: avatarOID})
	}

	_, err := collection.UpdateByID(ctx, oID, bson.D{{Key: "$set", Value: update}})

	if err != nil {
		return utils.ErrCannotUpdateUser
	}

	return nil

}

func (s *store) DeleteRefreshToken(ctx context.Context, refresh_token string) error {
	collection := s.db.Database(DbName).Collection(CollNameRefreshToken)

	_, err := collection.DeleteOne(ctx, bson.M{"token": refresh_token})

	if err != nil {
		return utils.ErrCannotLogout
	}

	return nil
}

func (s *store) CreateRefreshToken(ctx context.Context, p *models.RefreshTokenCreation) error {
	collection := s.db.Database(DbName).Collection(CollNameRefreshToken)

	userID, _ := primitive.ObjectIDFromHex(p.UserId)

	data := &models.RefreshToken{
		Token:  p.Token,
		UserId: userID,
		Exp:    p.Exp,
	}

	_, err := collection.InsertOne(ctx, data)

	if err != nil {
		return utils.ErrCannotCreateEntity
	}

	return nil
}

func (s *store) CheckUserByEmail(ctx context.Context, email string) (bool, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *store) GetUserByEmail(ctx context.Context, email string) (*models.UserModel, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "email", Value: email}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "avatar_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "avatar"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$avatar"},
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
		var item models.UserModel
		if err := cursor.Decode(&item); err != nil {
			fmt.Println(err)

			return nil, utils.ErrFailedToDecode
		}

		return &item, nil
	}

	if err := cursor.Err(); err != nil {
		return nil, utils.ErrCursorError
	}

	return nil, nil

}

func (s *store) GetUserByID(ctx context.Context, id string) (*models.UserModel, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	oid, _ := primitive.ObjectIDFromHex(id)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: oid}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "medias"},
					{Key: "localField", Value: "avatar_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "avatar"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$avatar"},
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
		var item models.UserModel
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

func (s *store) CreateUser(ctx context.Context, user *models.UserCreation) (primitive.ObjectID, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	data := &models.User{
		Email:    user.Email,
		Password: user.Password,
		Username: user.Username,
	}

	newUser, err := collection.InsertOne(ctx, data)

	if err != nil {
		return primitive.NilObjectID, utils.ErrCannotCreateEntity
	}

	id := newUser.InsertedID.(primitive.ObjectID)

	return id, nil
}
