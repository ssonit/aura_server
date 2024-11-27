package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	Media_Model "github.com/ssonit/aura_server/internal/media/models"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Username  string             `json:"username" bson:"username"`
	Bio       string             `json:"bio" bson:"bio"`
	AvatarID  primitive.ObjectID `json:"avatar_id" bson:"avatar_id"`
	Website   string             `json:"website" bson:"website"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Username  string             `json:"username" bson:"username"`
	Bio       string             `json:"bio" bson:"bio"`
	AvatarID  primitive.ObjectID `json:"avatar_id" bson:"avatar_id"`
	Avatar    *Media_Model.Media `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Website   string             `json:"website" bson:"website"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	BannedAt  *time.Time         `json:"banned_at,omitempty" bson:"banned_at,omitempty"`
}

func (m *User) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()

	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	if m.AvatarID.IsZero() {
		var err error
		m.AvatarID, err = primitive.ObjectIDFromHex("66e44a6c9cdaffa00a97a3dc")
		if err != nil {
			return nil, err
		}
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

type UserUpdate struct {
	AvatarID string `json:"avatar_id" bson:"avatar_id"`
	Username string `json:"username" bson:"username"`
	Bio      string `json:"bio" bson:"bio"`
	Website  string `json:"website" bson:"website"`
}

type RefreshToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Token     string             `json:"token" bson:"token"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Exp       time.Time          `json:"exp" bson:"exp"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (m *RefreshToken) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my RefreshToken
	return bson.Marshal((*my)(m))
}

type RefreshTokenCreation struct {
	Token  string    `json:"token" bson:"token"`
	UserId string    `json:"user_id" bson:"user_id"`
	Exp    time.Time `json:"exp" bson:"exp"`
}

type RefreshTokenSelection struct {
	Token string `json:"token" bson:"token"`
}
