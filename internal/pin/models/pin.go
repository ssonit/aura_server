package models

import (
	"net/url"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	User_Model "github.com/ssonit/aura_server/internal/auth/models"
	Board_Model "github.com/ssonit/aura_server/internal/board/models"
	Media_Model "github.com/ssonit/aura_server/internal/media/models"
)

type Pin struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
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

type PinModel struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	UserId      primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Title       string               `json:"title" bson:"title"`
	Description string               `json:"description" bson:"description"`
	MediaId     primitive.ObjectID   `json:"media_id" bson:"media_id"`
	LinkUrl     string               `json:"link_url" bson:"link_url"`
	User        User_Model.UserModel `json:"user,omitempty" bson:"user,omitempty"`
	Media       Media_Model.Media    `json:"media,omitempty" bson:"media,omitempty"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}

type PinCreation struct {
	BoardId     primitive.ObjectID `json:"board_id" bson:"board_id"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	MediaId     primitive.ObjectID `json:"media_id" bson:"media_id"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
}

type PinUpdate struct {
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
	BoardId     primitive.ObjectID `json:"board_id" bson:"board_id"`
}

type Filter struct {
	Title   string `json:"title,omitempty" bson:"title,omitempty" form:"title"`
	UserId  string `json:"user_id" bson:"user_id" form:"user_id"`
	SortKey string `json:"sort_key,omitempty" bson:"sort_key,omitempty" form:"sort_key"`
	Sort    string `json:"sort,omitempty" bson:"sort,omitempty" form:"sort"`
}

func (f *Filter) DecodeQuery() error {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.CanSet() {
			decodedValue, err := url.QueryUnescape(field.String())
			if err != nil {
				return err
			}
			field.SetString(decodedValue)
		}
	}
	return nil
}

type BoardPin struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	BoardId   primitive.ObjectID `json:"board_id" bson:"board_id"`
	PinId     primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

func (m *BoardPin) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.ID.IsZero() {
		m.ID = primitive.NewObjectID()
	}

	type my BoardPin
	return bson.Marshal((*my)(m))
}

type BoardPinModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	BoardId   primitive.ObjectID `json:"board_id" bson:"board_id"`
	PinId     primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Board     Board_Model.Board  `json:"board,omitempty" bson:"board,omitempty"`
	Pin       Pin                `json:"pin,omitempty" bson:"pin,omitempty"`
	Media     Media_Model.Media  `json:"media,omitempty" bson:"media,omitempty"`
	User      User_Model.User    `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type BoardPinCreation struct {
	BoardId primitive.ObjectID `json:"board_id" bson:"board_id"`
	PinId   primitive.ObjectID `json:"pin_id" bson:"pin_id"`
	UserId  primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type BoardPinFilter struct {
	BoardId primitive.ObjectID `json:"board_id,omitempty" bson:"board_id,omitempty"`
	PinId   primitive.ObjectID `json:"pin_id,omitempty" bson:"pin_id,omitempty"`
	UserId  primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type BoardPinUpdate struct {
	BoardId primitive.ObjectID `json:"board_id" bson:"board_id"`
}
