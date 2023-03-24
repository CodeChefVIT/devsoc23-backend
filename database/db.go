package database

import (
	"context"
	controller "devsoc23-backend/controllers"

	"os"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() controller.Database {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		panic(err.Error())
	}
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_ADDR"),
		Password: os.Getenv("REDIS_DB_PASS"),
		DB:       1,
	})
	return controller.Database{
		MongoClient: client,
		RedisClient: rdb,
	}

}