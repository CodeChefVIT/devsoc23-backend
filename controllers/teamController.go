package controller

import (
	"context"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (databaseClient Database) CreateTeam(ctx *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Id does not exist"})
	}

	id, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	findUser := models.User{}
	filter := bson.M{"_id": id}

	if err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	var payload *models.CreateTeamRequest

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	now := time.Now()

	newTeam := models.Team{
		Id:               primitive.NewObjectID(),
		TeamName:         payload.TeamName,
		TeamLeaderId:     findUser.Id,
		TeamMembers:      payload.TeamMembers,
		TeamSize:         payload.TeamSize,
		ProjectId:        primitive.NewObjectID(),
		InvitedTeammates: payload.InvitedTeammates,
		Round:            payload.Round,
		IsFinalised:      false,
		InviteLink:       payload.InviteLink,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	result, err := teamCollection.InsertOne(context.TODO(), newTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println(result.InsertedID)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": result})

}

func (databaseClient Database) GetTeam(ctx *fiber.Ctx) error {

	var teamName string

	if err := ctx.BodyParser(&teamName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filter := bson.M{"teamName": teamName}

	err := teamCollection.FindOne(context.TODO(), filter).Decode(&findTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}
	fmt.Println(findTeam.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": findTeam})
}

func (databaseClient Database) GetTeams(ctx *fiber.Ctx) error {

	var teamCollection = databaseClient.MongoClient.Database("devsoc").Collection("teams")
	var teams []models.Team
	cur, err := teamCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()

		var team models.Team
		err := cur.Decode(&team)
		if err != nil {
			log.Fatal(err)
		}

		teams = append(teams, team)

	}
	if err := cur.Err(); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "users": teams})
}

func (databaseClient Database) GetTeamMembers(ctx *fiber.Ctx) error {

	var teamName string

	if err := ctx.BodyParser(&teamName); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filter := bson.M{"TeamName": teamName}

	err := teamCollection.FindOne(context.TODO(), filter).Decode(&findTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}
	fmt.Println(findTeam.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "teamMembers": findTeam.TeamMembers})
}

/*
func (databaseClient Database) UpdateTeam(ctx *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Id does not exist"})
	}

	id, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	findUser := models.User{}
	filter := bson.M{"_id": id}

	if err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	update := bson.M{{"$set"}, bson.M{{}}}

	var payload *models.CreateTeamRequest

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	now := time.Now()

	if findUser.Id != id {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not TeamLeader"})
	}

	result, err := teamCollection.UpdateOne(context.TODO(), newTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println(result.InsertedID)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": result})

}
*/
/*

func (databaseClient Database) DeleteTeam(ctx *fiber.Ctx) error {

	return
}

func (databaseClient Database) JoinTeam(ctx *fiber.Ctx) error {

	return
}

func (databaseClient Database) FinaliseTeam(ctx *fiber.Ctx) error {

	return
}
*/
