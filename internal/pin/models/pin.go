package models

import (
	"net/url"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	User_Model "github.com/ssonit/aura_server/internal/auth/models"
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
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	MediaId     primitive.ObjectID `json:"media_id" bson:"media_id"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
	User        User_Model.User    `json:"user,omitempty" bson:"user,omitempty"`
	Media       Media_Model.Media  `json:"media,omitempty" bson:"media,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type PinCreation struct {
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	MediaId     primitive.ObjectID `json:"media_id" bson:"media_id"`
	LinkUrl     string             `json:"link_url" bson:"link_url"`
}

type Filter struct {
	Title string `json:"title,omitempty" bson:"title,omitempty" form:"title"`
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
