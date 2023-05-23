package database

import (
	"context"
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/initializers"

	"os"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() controller.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:password@mongo:27017/devsoc?authSource=admin"))
	if err != nil {
		panic(err.Error())
	}
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Username: "default", // use your Redis user. More info https://redis.io/docs/management/security/acl/
		Addr:     os.Getenv("REDIS_DB_ADDR"),
		Password: os.Getenv("REDIS_DB_PASS"),
		DB:       0,
	})
	s3Client := initializers.InitializeSpaces()
	return controller.Database{
		MongoClient: client,
		RedisClient: rdb,
		S3Client:    s3Client,
	}

}
