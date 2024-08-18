package common

import (
	"time"
)

type CommonModel struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
