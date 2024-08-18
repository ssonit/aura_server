package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Username  string    `json:"username" bson:"username"`
	Bio       string    `json:"bio" bson:"bio"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	Website   string    `json:"website" bson:"website"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UserCreation struct {
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Username  string    `json:"username" bson:"username"`
	Bio       string    `json:"bio" bson:"bio"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	Website   string    `json:"website" bson:"website"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *UserCreation) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	if u.Bio == "" {
		u.Bio = ""
	}
	if u.Avatar == "" {
		u.Avatar = ""
	}
	if u.Website == "" {
		u.Website = ""
	}

	type my UserCreation
	return bson.Marshal((*my)(u))
}
