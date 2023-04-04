package controller

import (
	"context"
	"devsoc23-backend/helper"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"log"
	"reflect"
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

	LeaderId, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	Leader := models.User{}
	filter := bson.M{"_id": LeaderId}

	err = userCollection.FindOne(context.TODO(), filter).Decode(&Leader)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	if Leader.InTeam {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Already in team"})
	}

	var payload models.CreateTeamRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": errors})
	}
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	now := time.Now()
	round := 1
	inviteCode := helper.RandSeq(6)
	IdArr := []primitive.ObjectID{LeaderId}

	newTeam := models.Team{
		Id:            primitive.NewObjectID(),
		TeamName:      payload.TeamName,
		TeamLeaderId:  LeaderId,
		TeamMembers:   IdArr,
		TeamSize:      1,
		ProjectExists: false,
		Round:         round,
		IsFinalised:   false,
		InviteCode:    inviteCode,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	result, err := teamCollection.InsertOne(context.TODO(), newTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team Not Inserted"})
	}

	filterUser := bson.M{"_id": LeaderId}
	updateUser := bson.M{"$set": bson.M{"inteam": true, "teamid": newTeam.Id}}

	userResult, err := userCollection.UpdateOne(context.TODO(), filterUser, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": result, "user": userResult})

}

func (databaseClient Database) GetTeam(ctx *fiber.Ctx) error {

	Id := ctx.Params("teamId")
	teamId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filter := bson.M{"_id": teamId}

	err = teamCollection.FindOne(context.TODO(), filter).Decode(&findTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": findTeam})
}

func (databaseClient Database) GetTeams(ctx *fiber.Ctx) error {

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

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
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "teams": teams})

}

func (databaseClient Database) GetTeamMembers(ctx *fiber.Ctx) error {

	id := ctx.Params("teamId")
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filter := bson.M{"_id": Id}

	err = teamCollection.FindOne(context.TODO(), filter).Decode(&findTeam)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "teammembers": findTeam.TeamMembers})
}

func (databaseClient Database) UpdateTeam(ctx *fiber.Ctx) error {

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No Id passed"})
	}

	userId, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	Id := ctx.Params("teamId")
	teamId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}

	oldTeam := models.Team{}
	filterTeamId := bson.M{"_id": teamId}

	if err := teamCollection.FindOne(context.TODO(), filterTeamId).Decode(&oldTeam); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId incorrect"})
	}

	if oldTeam.TeamLeaderId != userId {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "user not teamleader"})
	}

	var payload *models.UpdateTeam

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": errors})
	}

	now := time.Now()

	filterTeamUpdate := bson.M{"_id": teamId}

	// var setUpdates bson.D

	var setUpdates bson.D

	setUpdates = append(setUpdates,
		bson.E{"updatedat", now})

	if payload.TeamName != nil {
		setUpdates = append(setUpdates, bson.E{"teamname", payload.TeamName})
	}
	if payload.IsFinalised != oldTeam.IsFinalised {
		setUpdates = append(setUpdates, bson.E{"isfinalised", payload.IsFinalised})
	}
	if payload.ProjectExists != oldTeam.ProjectExists {
		setUpdates = append(setUpdates, bson.E{"projectexists", payload.ProjectExists})
	}
	if payload.Round != oldTeam.Round {
		setUpdates = append(setUpdates, bson.E{"round", payload.Round})
	}
	if len(payload.InviteCode) > 0 {
		setUpdates = append(setUpdates, bson.E{"inviteCode", payload.InviteCode})
	}

	result, err := teamCollection.UpdateOne(context.TODO(), filterTeamUpdate, bson.D{{"$set", setUpdates}})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}
	fmt.Println(payload.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "team": result})
}

func (databaseClient Database) DeleteTeam(ctx *fiber.Ctx) error {
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No Id passed"})
	}

	userId, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	Id := ctx.Params("teamId")
	teamId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}

	delTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	if err := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&delTeam); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId incorrect"})
	}

	if delTeam.TeamLeaderId != userId {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "user not teamleader"})
	}

	filterUser := bson.M{"teamid": delTeam.Id}

	updateUserInTeam := bson.M{"$set": bson.M{"inteam": false}}

	_, err = userCollection.UpdateMany(context.TODO(), filterUser, updateUserInTeam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "users not found"})
	}

	delFilter := bson.M{"_id": delTeam.Id}
	result, err := teamCollection.DeleteOne(context.TODO(), delFilter)
	fmt.Println("DeleteOne Result:", reflect.TypeOf(result))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not deleted"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Team deleted"})
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
	Id := ctx.Params("teamId")
	teamId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}

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
		{Key: "$push", Value: bson.D{{Key: "teamMembers", Value: id}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		log.Fatal(err)
	}
	update := bson.M{"inteam": true, "teamid": teamId}
	updateUser := bson.M{
		"$set": update,
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User Update failed"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": result, "team message": res})
}

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
	Id := findUser.TeamId
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	teamId, err := primitive.ObjectIDFromHex(*Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}

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
		{Key: "$pull", Value: bson.D{{Key: "teamMembers", Value: id}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		log.Fatal(err)
	}
	update := bson.M{"inteam": false, "teamid": ""}
	updateUser := bson.M{
		"$set": update,
	}

	result, err := userCollection.UpdateOne(context.TODO(), filter, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User leaving failed"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": result, "team message": res})
}

func (databaseClient Database) PromoteTeam(ctx *fiber.Ctx) error {

	teamId := ctx.Params("teamId")
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	errr := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&findTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}

	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "round", Value: findTeam.Round + 1}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": errr})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": res})
}

func (databaseClient Database) FinaliseTeam(ctx *fiber.Ctx) error {

	teamId := ctx.Params("teamId")
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	errr := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&findTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}

	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "isFinalised", Value: true}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": errr})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": res})
}

func (databaseClient Database) DisqualifyTeam(ctx *fiber.Ctx) error {

	teamId := ctx.Params("teamId")
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	findTeam := models.Team{}
	filterTeam := bson.M{"_id": teamId}

	errr := teamCollection.FindOne(context.TODO(), filterTeam).Decode(&findTeam)

	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not found"})
	}

	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "round", Value: 0}}},
	}

	res, errr := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": errr})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": res})
}
