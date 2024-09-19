package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	PinId     primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (m *Like) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my Like
	return bson.Marshal((*my)(m))
}

type LikeCreation struct {
	UserId string `json:"user_id" bson:"user_id"`
	PinId  string `json:"pin_id" bson:"pin_id"`
}

type LikeDelete struct {
	UserId string `json:"user_id" bson:"user_id"`
	PinId  string `json:"pin_id" bson:"pin_id"`
}
