package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Username  string    `json:"username" bson:"username"`
	Bio       string    `json:"bio,omitempty" bson:"bio,omitempty"`
	Avatar    string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Website   string    `json:"website,omitempty" bson:"website,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UserCreation struct {
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Username  string    `json:"username" bson:"username"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (u *UserCreation) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	type my UserCreation
	return bson.Marshal((*my)(u))
}
