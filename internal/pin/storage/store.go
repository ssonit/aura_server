package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "aura_pins"
	CollName = "pin"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}
