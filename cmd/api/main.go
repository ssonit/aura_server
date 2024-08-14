package main

import (
	"context"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ssonit/aura_server/common"
	"github.com/ssonit/aura_server/middleware"

	"go.uber.org/zap"

	pin_http "github.com/ssonit/aura_server/internal/pin/transport/gin"
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
	mongoAddr = common.EnvConfig("MONGO_DB_HOST", "localhost:27017")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	// uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoAddr)
	// fmt.Println(uri)
	// mongoClient, err := connectMongoDB(uri)
	// if err != nil {
	// 	logger.Fatal("failed to connect to mongo db", zap.Error(err))
	// }

	// fmt.Println(mongoClient)

	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.Recovery())

	h := pin_http.NewHandler()
	h.RegisterRoutes(r)

	logger.Info("Server listening on ", zap.String("port", httpAddr))

	r.Run(httpAddr)
}
