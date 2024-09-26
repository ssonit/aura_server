package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suggestion struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Count     int                `json:"count" bson:"count"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
