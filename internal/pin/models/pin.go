package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pin struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	MediaId     primitive.ObjectID `json:"media_id" bson:"media_id"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *Pin) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my Pin
	return bson.Marshal((*my)(m))
}

type PinCreation struct {
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	MediaId     primitive.ObjectID `json:"media_id" bson:"media_id"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
}
