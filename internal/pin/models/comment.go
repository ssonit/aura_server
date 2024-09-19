package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	User_Model "github.com/ssonit/aura_server/internal/auth/models"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	PinId     primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *Comment) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my Comment
	return bson.Marshal((*my)(m))
}

type CommentModel struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	UserId    primitive.ObjectID   `json:"user_id" bson:"user_id"`
	PinId     primitive.ObjectID   `json:"pin_id" bson:"pin_id"`
	Content   string               `json:"content" bson:"content"`
	User      User_Model.UserModel `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

type CommentCreation struct {
	PinId   primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	UserId  primitive.ObjectID `json:"user_id" bson:"user_id"`
	Content string             `json:"content" bson:"content"`
}
