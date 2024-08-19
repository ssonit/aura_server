package storage

import (
	"context"

	"github.com/ssonit/aura_server/internal/auth/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "aura_pins"
	CollName = "users"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db: db}
}

func (s *store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *store) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	collection := s.db.Database(DbName).Collection(CollName)

	var user models.User

	oid, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	if err != nil {
		return nil, err
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

	id := newUser.InsertedID.(primitive.ObjectID)

	return id, err
}
