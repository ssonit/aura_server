package storage

import (
	"context"

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

func (s *store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	return &user, nil
}

func (s *store) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	var user models.User

	oid, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	return &user, nil
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
