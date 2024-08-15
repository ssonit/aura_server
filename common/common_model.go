package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommonModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IsDeleted bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
