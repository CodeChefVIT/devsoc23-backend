package controller

import (
	"context"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (databaseClient Database) CreateTimeLine(ctx *fiber.Ctx) error {
	var payload *models.TimeLine

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Timeline Collection
	timelineColl := databaseClient.MongoClient.Database("devsoc").Collection("timelines")

	newTimeline := models.TimeLine{
		Id:          primitive.NewObjectID(),
		Title:       payload.Title,
		SubTitle:    payload.SubTitle,
		Description: payload.Description,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
	}

	_, err := timelineColl.InsertOne(context.TODO(), newTimeline)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Timeline created succesfully", "Timeline": newTimeline})
}

func (databaseClient Database) GetAllTimeLine(ctx *fiber.Ctx) error {
	var timelineColl = databaseClient.MongoClient.Database("devsoc").Collection("timelines")
	var timelines []models.TimeLine
	cur, err := timelineColl.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()

		var timeline models.TimeLine
		err := cur.Decode(&timeline)
		if err != nil {
			log.Fatal(err)
		}

		timelines = append(timelines, timeline)

	}
	if err := cur.Err(); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "timelines": timelines})
}

func (databaseClient Database) GetTimeLine(ctx *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}
	timelineColl := databaseClient.MongoClient.Database("devsoc").Collection("timelines")

	filter := bson.M{"_id": id}
	findTimeline := models.TimeLine{}

	err = timelineColl.FindOne(context.TODO(), filter).Decode(&findTimeline)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Timeline not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "Timeline Found": findTimeline})
}

func (databaseClient Database) UpdateTimeLine(ctx *fiber.Ctx) error {

	// Get request body and bind to payload
	var payload *models.TimeLine
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Get timeline id from params
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	// Create update map
	update := bson.M{
		"title":       payload.Title,
		"subtitle":    payload.SubTitle,
		"description": payload.Description,
		"starttime":   payload.StartTime,
		"endtime":     payload.EndTime,
	}

	timelineColl := databaseClient.MongoClient.Database("devsoc").Collection("timelines")

	// Update timeline in timeline document
	_, err = timelineColl.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update timeline"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Timeline updated succesfully"})
}

func (databaseClient Database) DeleteTimeline(ctx *fiber.Ctx) error {

	// Get timeline id from params
	id, err := primitive.ObjectIDFromHex(ctx.Params("id"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	timelineColl := databaseClient.MongoClient.Database("devsoc").Collection("timelines")

	// Delete user from user document
	filter := bson.M{"_id": id}
	deleteResult, _ := timelineColl.DeleteOne(context.TODO(), filter)

	if deleteResult.DeletedCount == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Timeline not deleted"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Timeline deleted succesfully"})
}
