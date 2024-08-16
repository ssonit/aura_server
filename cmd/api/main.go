package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/internal/server"

	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	return client, err
}

var (
	httpAddr  = common.EnvConfig("HTTP_ADDR", ":8080")
	mongoUser = common.EnvConfig("MONGO_DB_USERNAME", "root")
	mongoPass = common.EnvConfig("MONGO_DB_PASSWORD", "admin")
	mongoAddr = common.EnvConfig("MONGO_DB_HOST", "mongodb:27017")
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoAddr)
	fmt.Println(uri)
	mongoClient, err := connectMongoDB(uri)
	if err != nil {
		logger.Fatal("failed to connect to mongodb", zap.Error(err))
	}

	r := gin.Default()

	s := server.NewServer(r, mongoClient, logger)
	if err = s.Run(httpAddr); err != nil {
		log.Fatal(err)
	}

}
