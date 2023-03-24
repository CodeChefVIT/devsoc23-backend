package controller

import (
	"context"
	"devsoc23-backend/helper"
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
	inviteCode := helper.RandSeq(6)

	newTeam := models.Team{
		Id:               primitive.NewObjectID(),
		TeamName:         payload.TeamName,
		TeamLeaderId:     findUser.Id,
		TeamMembers:      payload.TeamMembers,
		TeamSize:         1,
		ProjectId:        primitive.NewObjectID(),
		InvitedTeammates: payload.InvitedTeammates,
		Round:            payload.Round,
		IsFinalised:      false,
		InviteCode:       inviteCode,
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

func (databaseClient Database) JoinTeam(ctx *fiber.Ctx) error {

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

	if findUser.InTeam {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "you are already in team"})
	}
	teamId := ctx.Params("teamId")
	inviteCode := ctx.Params("inviteCode")
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	errr := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&findTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}

	if id == findTeam.TeamLeaderId {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "you cannot invite yourself for your team"})
	}

	if findTeam.InviteCode != inviteCode {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong Invite Code"})

	}
	if findTeam.TeamSize == 4 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team is full"})
	}

	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "teamSize", Value: findTeam.TeamSize + 1}}},
		{Key: "$push", Value: bson.D{{Key: "teamMember", Value: id}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		log.Fatal(err)
	}
	update := bson.M{"inTeam": true, "teamId": teamId}
	updateUser := bson.M{
		"$set": update,
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User Update failed"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": result, "team message": res})
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

// func (databaseClient Database) DeleteTeam(ctx *fiber.Ctx) error {

// 	return
// }

// func (databaseClient Database) FinaliseTeam(ctx *fiber.Ctx) error {

// 	return
// }

func (databaseClient Database) LeaveTeam(ctx *fiber.Ctx) error {

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

	if !findUser.InTeam {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "you are in no team"})
	}
	teamId := findUser.TeamId
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	errr := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&findTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}

	if id == findTeam.TeamLeaderId {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "you cannot leave your team"})
	}


	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "teamSize", Value: findTeam.TeamSize - 1}}},
		{Key: "$pull", Value: bson.D{{Key: "teamMember", Value: id}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		log.Fatal(err)
	}
	update := bson.M{"inTeam": false, "teamId": ""}
	updateUser := bson.M{
		"$set": update,
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User leaving failed"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": result, "team message": res})
}