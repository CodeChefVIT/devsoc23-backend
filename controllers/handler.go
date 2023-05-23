package controller

import (
	"cloud.google.com/go/storage"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	MongoClient *mongo.Client
	RedisClient *redis.Client
	S3Client    *storage.Client
}
