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
	Type      string             `json:"type" bson:"type"` // "all_pins" hoặc "custom"
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // Thêm DeletedAt để soft delete
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
	Type      string             `json:"type" bson:"type"` // "all_pins" hoặc "custom"
}

type Filter struct {
	UserId    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	IsPrivate bool               `json:"isPrivate,omitempty" bson:"isPrivate,omitempty"`
}

type BoardModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	IsPrivate bool               `json:"isPrivate" bson:"isPrivate"`
	Type      string             `json:"type" bson:"type"` // "all_pins" hoặc "custom"
	User      User_Model.User    `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type BoardUpdate struct {
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	IsPrivate bool   `json:"isPrivate,omitempty" bson:"isPrivate,omitempty"`
}
