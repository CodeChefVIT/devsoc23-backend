package controller

import (
	"context"
	"devsoc23-backend/helper"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"log"
	"os"
	"strings"
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

	filter = bson.M{"teamname": payload.TeamName}
	count, _ := teamCollection.CountDocuments(context.TODO(), filter)
	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "Name already exists"})
	}

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

func (databaseClient Database) GetIsMember(ctx *fiber.Ctx) error {

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	token := fields[1]

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	filter := bson.M{"_id": res.Id}

	findUser := models.User{}
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	err = userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
	}

	if !findUser.InTeam {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "inTeam": false})
	}
	teamId, err := primitive.ObjectIDFromHex(*findUser.TeamId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}
	var users []models.User
	cur, err := userCollection.Find(context.TODO(), bson.M{"teamid": teamId})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {

		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "message": "User not decodable"})
		}

		users = append(users, user)

	}
	if err := cur.Err(); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "inTeam": true, "memberDetails": users})
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
		bson.E{Key: "updatedat", Value: now})

	if payload.TeamName != nil {
		filter := bson.M{"teamname": payload.TeamName}
		count, _ := teamCollection.CountDocuments(context.TODO(), filter)
		if count > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "Name already exists"})
		}
		setUpdates = append(setUpdates, bson.E{Key: "teamname", Value: payload.TeamName})
	}
	if payload.IsFinalised != oldTeam.IsFinalised {
		setUpdates = append(setUpdates, bson.E{Key: "isfinalised", Value: payload.IsFinalised})
	}
	if payload.ProjectExists != oldTeam.ProjectExists {
		setUpdates = append(setUpdates, bson.E{Key: "projectexists", Value: payload.ProjectExists})
	}
	if payload.Round != oldTeam.Round {
		setUpdates = append(setUpdates, bson.E{Key: "round", Value: payload.Round})
	}
	if len(payload.InviteCode) > 0 {
		setUpdates = append(setUpdates, bson.E{Key: "inviteCode", Value: payload.InviteCode})
	}

	result, err := teamCollection.UpdateOne(context.TODO(), filterTeamUpdate, bson.D{{Key: "$set", Value: setUpdates}})

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error()})
	}

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

	updateUserInTeam := bson.M{"$set": bson.M{"inteam": false, "teamid": ""}}

	_, err = userCollection.UpdateMany(context.TODO(), filterUser, updateUserInTeam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "users not found"})
	}

	delFilter := bson.M{"_id": delTeam.Id}
	result, err := teamCollection.DeleteOne(context.TODO(), delFilter)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not deleted"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": result})
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

	// if id == findTeam.TeamLeaderId {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "teamleader cannot leave team"})
	// }

	if id == findTeam.TeamLeaderId {
		// leader trying to remove himself
		// delete the team itself
		// return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Leader cannot exit the team"})
		delFilter := bson.M{"_id": teamId}
		_, err := teamCollection.DeleteOne(context.TODO(), delFilter)

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not deleted"})
		}

		filterUser := bson.M{"teamid": teamId}
		updateUserInTeam := bson.M{"$set": bson.M{"inteam": false, "teamid": ""}}
		// userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
		_, err = userCollection.UpdateMany(context.TODO(), filterUser, updateUserInTeam)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "users not found"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "TeamLeader left, Team deleted"})
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

	filter = bson.M{"_id": id}

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

func (databaseClient Database) RemoveMember(ctx *fiber.Ctx) error {

	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No Id passed"})
	}

	leaderId, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}

	Id := ctx.Params("teamId")
	teamId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId not parsable"})
	}

	Team := models.Team{}
	filterTeamId := bson.M{"_id": teamId}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")

	if err := teamCollection.FindOne(context.TODO(), filterTeamId).Decode(&Team); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "TeamId incorrect"})
	}

	if Team.TeamLeaderId != leaderId {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "user not teamleader"})
	}

	memberId := ctx.Params("memberId")
	remMemberId, err := primitive.ObjectIDFromHex(memberId)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Member Id not parsable"})
	}

	if remMemberId == leaderId {
		// leader trying to remove himself
		// delete the team itself
		// return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Leader cannot exit the team"})
		delFilter := bson.M{"_id": Team.Id}
		_, err := teamCollection.DeleteOne(context.TODO(), delFilter)

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Team not deleted"})
		}

		filterUser := bson.M{"teamid": Team.Id}
		updateUserInTeam := bson.M{"$set": bson.M{"inteam": false, "teamid": ""}}
		userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
		_, err = userCollection.UpdateMany(context.TODO(), filterUser, updateUserInTeam)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "users not found"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Team deleted"})

	}

	var newTeamSize int
	var newTeamMembers []primitive.ObjectID

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	remUser := models.User{}
	filter := bson.M{"_id": remMemberId}

	if err := userCollection.FindOne(context.TODO(), filter).Decode(&remUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	if !remUser.InTeam {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "member is not in any team"})
	}

	present := false
	for index, id := range Team.TeamMembers {
		if id == remMemberId {
			present = true
		}
		if present {
			newTeamMembers = append(Team.TeamMembers[:index], Team.TeamMembers[index+1:]...)
			newTeamSize = Team.TeamSize - 1
			break
		}
	}

	if !present {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Member not present in given team"})
	}

	updateTeam := bson.D{
		{Key: "$set", Value: bson.D{{Key: "teamSize", Value: newTeamSize}}},
		{Key: "$set", Value: bson.D{{Key: "teamMembers", Value: newTeamMembers}}},
	}

	res, err := teamCollection.UpdateOne(context.TODO(), bson.M{"_id": teamId}, updateTeam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Error in updating fields"})
	}

	update := bson.M{"inteam": false, "teamid": ""}
	updateUser := bson.M{
		"$set": update,
	}

	result, err := userCollection.UpdateOne(context.TODO(), bson.M{"_id": remMemberId}, updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User leaving failed"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Member removed", "result1": res, "result2": result})

}
