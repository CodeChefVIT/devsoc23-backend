package database

import (
	"context"
	controller "devsoc23-backend/controllers"
	// "fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
)

// Database struct

func NewDatabase() controller.Database {

	// USER := os.Getenv("DB_USER")
	// PASS := os.Getenv("DB_PASSWORD")
	// HOST := os.Getenv("DB_HOST")
	// DBNAME := os.Getenv("DB_NAME")
	// PORT := os.Getenv("DB_PORT")

	// URL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", HOST, USER, PASS,
	// 	DBNAME, PORT)
	// db, err := gorm.Open(postgres.Open(URL), &gorm.Config{})

	// if err != nil {
	// 	panic("Failed to connect to database!")

	// }
	// fmt.Println("Database connection established")

	credential := options.Credential{
		Username: "root",
		Password: "password",
	}
	_ = credential
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(credential))
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
		// DB:          db,
		MongoClient: client,
		RedisClient: rdb,
	}
}