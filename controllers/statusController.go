package controller

import (
	"context"
	"devsoc23-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (databaseClient Database) CheckIn(ctx *fiber.Ctx) error {

	statusCollection := databaseClient.MongoClient.Database("devsoc").Collection("status")
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("user")

	userId := ctx.Params("userId")
	userid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		panic(err)
	}
	findUser := models.Status{}

	errr := statusCollection.FindOne(context.TODO(), userId).Decode(&findUser)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User status not found"})
	}
	if !findUser.InHall {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "user is already checkedout"})
	}
	teamid := findUser.TeamId
	now := time.Now()
	t := now.Format("2006-01-02 15:04:05")
	index := 0
	status := "checkIn"

	newStatus := models.Status{
		Id:        primitive.NewObjectID(),
		UserId:    userid,
		TeamId:    teamid,
		InHall:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	newStatus.Time.Num = index
	newStatus.Time.IsTime = t
	newStatus.Time.IfStatus = status

	statusRes, err := statusCollection.UpdateByID(context.TODO(), userid, newStatus)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	userRes, err := userCollection.UpdateByID(context.TODO(), userid, bson.M{"$set": bson.M{"IsCheckedIn": true}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "userStatus": statusRes, "user": userRes})
}

func (databaseClient Database) CheckOut(ctx *fiber.Ctx) error {

	statusCollection := databaseClient.MongoClient.Database("devsoc").Collection("status")
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("user")

	userId := ctx.Params("userId")
	userid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		panic(err)
	}
	findUser := models.Status{}

	errr := statusCollection.FindOne(context.TODO(), userid).Decode(&findUser)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User status not found"})
	}
	if !findUser.InHall {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "user is already checkedout"})
	}

	teamid := findUser.TeamId
	now := time.Now()
	t := now.Format("2006-01-02 15:04:05")
	status := "checkOut"
	index := 1

	newStatus := models.Status{
		Id:        primitive.NewObjectID(),
		UserId:    userid,
		TeamId:    teamid,
		InHall:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	newStatus.Time.Num = index
	newStatus.Time.IsTime = t
	newStatus.Time.IfStatus = status

	statusRes, err := userCollection.UpdateByID(context.TODO(), userid, newStatus)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "userStatus": statusRes})
}
