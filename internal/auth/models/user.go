package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Username  string             `json:"username" bson:"username"`
	Bio       string             `json:"bio" bson:"bio"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Website   string             `json:"website" bson:"website"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *User) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my User
	return bson.Marshal((*my)(m))
}

type UserLogin struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserCreation struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Username string `json:"username" bson:"username"`
}
