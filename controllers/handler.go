package controller

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Database struct {
	DB          *gorm.DB
	MongoClient *mongo.Client
	RedisClient *redis.Client
}
