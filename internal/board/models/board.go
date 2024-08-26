package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	User_Model "github.com/ssonit/aura_server/internal/auth/models"
)

type Board struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	IsPrivate bool               `json:"isPrivate" bson:"isPrivate"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *Board) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my Board
	return bson.Marshal((*my)(m))
}

type BoardCreation struct {
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	IsPrivate bool               `json:"isPrivate" bson:"isPrivate"`
}

type Filter struct {
	UserId primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type BoardModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	IsPrivate bool               `json:"isPrivate" bson:"isPrivate"`
	User      User_Model.User    `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
