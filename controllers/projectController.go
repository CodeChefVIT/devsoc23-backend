package controller

import (
	"context"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (databaseClient Database) GetTeamIdByUserId(ctx *fiber.Ctx) *string {
	idString := ctx.GetRespHeader("currentUser")
	//parsing userid
	if idString == "" {
		erms := "ID does not exist"
		return &erms
	}

	id, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		erms := "User id not parseable"
		return &erms
	}
	//getting user
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	findUser := models.User{}
	filter := bson.M{"_id": id}

	if err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser); err != nil {
		erms := "user not found"
		return &erms
	}
	//check if user is in a team:
	if !findUser.InTeam {
		erms := "user not in a team"
		return &erms
	}

	//getting team of the user
	team_id := findUser.TeamId

	return team_id
}

func (databaseClient Database) CreateProject(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}
	//validating request
	var payload *models.CreateProjectRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}
	//inserting project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	status := "Under Review"
	newProject := models.Project{
		Id:                 primitive.NewObjectID(),
		TeamId:             team_id,
		ProjectName:        payload.ProjectName,
		ProjectDescription: payload.ProjectDescription,
		ProjectStatus:      &status,
		ProjectVideoLink:   payload.ProjectVideoLink,
		ProjectGithubLink:  payload.ProjectGithubLink,
		ProjectTrack:       payload.ProjectTrack,
		ProjectTags:        payload.ProjectTags,
		IsFinal:            false,
	}

	result, err := projectCollection.InsertOne(context.TODO(), newProject)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	//updating project id in team:
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")
	team_filter := bson.M{"_id": team_id}
	team_update := bson.M{"$set": bson.M{"ProjectId": newProject.Id}}
	team_result, err := teamCollection.UpdateOne(context.Background(), team_filter, team_update)
	if err != nil {
		return err
	}
	if team_result.ModifiedCount == 0 {
		return fmt.Errorf("team not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": result})
}

func (databaseClient Database) GetProjectByUserid(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}

	//getting project
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	findProject := models.Project{}
	project_filter := bson.M{"TeamId": team_id}

	if err := projectCollection.FindOne(context.TODO(), project_filter).Decode(&findProject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No project found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": findProject})
}

func (databaseClient Database) GetProjectByTeamid(ctx *fiber.Ctx) error {

	teamId := ctx.Params("teamId")

	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")

	findProject := models.Project{}
	project_filter := bson.M{"TeamId": teamId}

	err := projectCollection.FindOne(context.TODO(), project_filter).Decode(&findProject)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Project not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": findProject})
}

func (databaseClient Database) UpdateProject(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}
	//validating request
	var payload *models.UpdateProjectRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	//updating project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	project_filter := bson.M{"TeamId": team_id}
	project_update := bson.M{"$set": bson.M{
		"ProjectName":        payload.ProjectName,
		"ProjectDescription": payload.ProjectDescription,
		"ProjectVideoLink":   payload.ProjectVideoLink,
		"ProjectGithubLink":  payload.ProjectGithubLink,
		"ProjectTrack":       payload.ProjectTrack,
		"ProjectTags":        payload.ProjectTags,
	}}
	project_result, err := projectCollection.UpdateOne(context.Background(), project_filter, project_update)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": project_result})
}

func (databaseClient Database) DeleteProject(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}

	//updating project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	project_filter := bson.M{"TeamId": team_id}
	res, err := projectCollection.DeleteOne(context.Background(), project_filter)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Unable to delete"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": res})
}

func (databaseClient Database) FinaliseProjectSubmission(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}

	//updating project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	project_filter := bson.M{"TeamId": team_id}
	project_update := bson.M{"$set": bson.M{"IsFinal": true}}
	project_result, err := projectCollection.UpdateOne(context.Background(), project_filter, project_update)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": project_result})
}

func (databaseClient Database) GetStatus(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	if len(*team_id) < 25 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": team_id})
	}

	//getting project
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	findProject := models.Project{}
	project_filter := bson.M{"TeamId": team_id}

	if err := projectCollection.FindOne(context.TODO(), project_filter).Decode(&findProject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No project found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project_status": findProject.ProjectStatus})
}

func (databaseClient Database) GetProjects(ctx *fiber.Ctx) error {
	var projectCollection = databaseClient.MongoClient.Database("devsoc").Collection("projects")
	var projects []models.Project
	cur, err := projectCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()

		var project models.Project
		err := cur.Decode(&project)
		if err != nil {
			log.Fatal(err)
		}

		projects = append(projects, project)

	}
	if err := cur.Err(); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "projects": projects})
}
