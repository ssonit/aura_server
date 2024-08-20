package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Url       string             `json:"url" bson:"url"`
	SecureUrl string             `json:"secure_url" bson:"secure_url"`
	PublicId  string             `json:"public_id" bson:"public_id"`
	Format    string             `json:"format" bson:"format"`
	Width     int                `json:"width" bson:"width"`
	Height    int                `json:"height" bson:"height"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *Media) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my Media
	return bson.Marshal((*my)(m))
}

type MediaCreation struct {
	Url       string `json:"url" bson:"url"`
	SecureUrl string `json:"secure_url" bson:"secure_url"`
	PublicId  string `json:"public_id" bson:"public_id"`
	Format    string `json:"format" bson:"format"`
	Width     int    `json:"width" bson:"width"`
	Height    int    `json:"height" bson:"height"`
}
