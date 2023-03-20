package controller

import (
	"context"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (databaseClient Database) CreateTeam(ctx *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	email := ctx.GetRespHeader("currentUser")

	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email does not exist"})
	}

	findUser := models.User{}
	filter := bson.M{"email": email}

	err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
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
		// TeamId:           helper.GenerateToken(),
	}

	result, err := teamCollection.InsertOne(context.TODO(), newTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println(result.InsertedID)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": result})

}

/*

func (databaseClient Database) GetTeam(ctx *fiber.Ctx) error {

	return
}

func (databaseClient Database) GetTeams(ctx *fiber.Ctx) error {

	return
}

func (databaseClient Database) UpdateTeam(ctx *fiber.Ctx) error {

	return
}
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
