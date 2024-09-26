package server

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	DbName              = "aura_pins"
	CollNameTags        = "tags"
	CollNameSuggestions = "suggestions"
	CollNamePins        = "pins"
	CollNameBoards      = "boards"
	CollNameComments    = "comments"
	CollNameBoardPin    = "board_pins"
	CollNameLikes       = "likes"
)

type Server struct {
	r      *gin.Engine
	db     *mongo.Client
	logger *zap.Logger
}

func NewServer(r *gin.Engine, db *mongo.Client, logger *zap.Logger) *Server {
	return &Server{
		r:      r,
		db:     db,
		logger: logger,
	}
}

func createLikesIndex(db *mongo.Client) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "pin_id", Value: -1}, {Key: "user_id", Value: -1}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameLikes).Indexes().CreateOne(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}
}

func createBoardPinIndex(db *mongo.Client) {
	indexModel := []mongo.IndexModel{
		{Keys: bson.D{{Key: "pin_id", Value: -1}, {Key: "user_id", Value: -1}}},
		{Keys: bson.D{{Key: "board_id", Value: -1}}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameBoardPin).Indexes().CreateMany(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}

}

func createCommentsIndex(db *mongo.Client) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "pin_id", Value: -1}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameComments).Indexes().CreateOne(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}

}

func createBoardsIndex(db *mongo.Client) {
	indexModel := []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: -1}}},
		{Keys: bson.D{{Key: "user_id", Value: -1}, {Key: "type", Value: -1}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameBoards).Indexes().CreateMany(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}

}

func createPinsIndex(db *mongo.Client) {
	indexModel := []mongo.IndexModel{
		{Keys: bson.D{{Key: "title", Value: -1}}},
		{Keys: bson.D{{Key: "description", Value: -1}}},
		{Keys: bson.D{{Key: "tags", Value: -1}}},
		{Keys: bson.D{{Key: "user_id", Value: -1}}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Database(DbName).Collection(CollNamePins).Indexes().CreateMany(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}

}

func createSuggestionsIndex(db *mongo.Client) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: -1}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameSuggestions).Indexes().CreateOne(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}

}

func createTagsIndex(db *mongo.Client) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: -1}},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.Database(DbName).Collection(CollNameTags).Indexes().CreateOne(
		ctx,
		indexModel,
	)

	if err != nil {
		panic(err)
	}
}

func (s *Server) Run(httpAddr string) error {
	// Create indexes
	createTagsIndex(s.db)
	createSuggestionsIndex(s.db)
	createBoardsIndex(s.db)
	createPinsIndex(s.db)
	createCommentsIndex(s.db)
	createLikesIndex(s.db)
	createBoardPinIndex(s.db)

	// Map routes
	if err := s.MapRoutes(s.r, httpAddr); err != nil {
		return err
	}

	return nil
}
