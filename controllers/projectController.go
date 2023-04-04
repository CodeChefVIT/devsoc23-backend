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
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
	}
	fmt.Println("team id is: ", *team_id)
	//validating request
	var payload *models.CreateProjectRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	//checking if project is there
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")
	team_objectid, err := primitive.ObjectIDFromHex(*team_id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	team_filter := bson.M{"_id": team_objectid}
	find_team := models.Team{}
	if err := teamCollection.FindOne(context.TODO(), team_filter).Decode(&find_team); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println("Projectid: ", find_team.ProjectId)
	if find_team.ProjectExists {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "project already exists"})
	}
	//inserting project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	status := "Under Review"
	count := 0
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
		LikeCount:          count,
		LikesId:            nil,
	}

	result, err := projectCollection.InsertOne(context.TODO(), newProject)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println("Inserted")

	//updating project id in team:
	team_update := bson.M{"$set": bson.M{"projectId": newProject.Id, "projectExists": true}}
	team_result, err := teamCollection.UpdateOne(context.Background(), team_filter, team_update)
	if err != nil {
		return err
	}
	fmt.Println(team_result)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": result})
}

func (databaseClient Database) GetProjectByUserid(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
	}
	fmt.Println("team id: ", *team_id)
	//getting project
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	findProject := models.Project{}
	project_filter := bson.M{"teamId": team_id}

	if err := projectCollection.FindOne(context.TODO(), project_filter).Decode(&findProject); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No project found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": findProject})
}

func (databaseClient Database) GetProjectByTeamid(ctx *fiber.Ctx) error {

	teamId := ctx.Params("teamId")

	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")

	findProject := models.Project{}
	project_filter := bson.M{"teamId": teamId}

	err := projectCollection.FindOne(context.TODO(), project_filter).Decode(&findProject)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Project not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": findProject})
}

func (databaseClient Database) UpdateProject(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
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
	project_filter := bson.M{"teamId": team_id}
	project_update := bson.M{"$set": bson.M{
		"projectname":        payload.ProjectName,
		"projectdescription": payload.ProjectDescription,
		"projectvideolink":   payload.ProjectVideoLink,
		"projectgithublink":  payload.ProjectGithubLink,
		"projecttrack":       payload.ProjectTrack,
		"projecttags":        payload.ProjectTags,
	}}
	project_result, err := projectCollection.UpdateOne(context.Background(), project_filter, project_update)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": project_result})
}

func (databaseClient Database) DeleteProject(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
	}

	//updating project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	project_filter := bson.M{"teamId": team_id}
	res, err := projectCollection.DeleteOne(context.Background(), project_filter)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Unable to delete"})
	}

	//updating team:
	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")
	team_objectid, err := primitive.ObjectIDFromHex(*team_id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	team_filter := bson.M{"_id": team_objectid}
	team_update := bson.M{"$unset": bson.M{"projectId": "1", "projectExists": false}}
	team_result, err := teamCollection.UpdateOne(context.Background(), team_filter, team_update)
	if err != nil {
		return err
	}
	fmt.Println(team_result)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": res})
}

func (databaseClient Database) FinaliseProjectSubmission(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
	}

	//updating project:
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	project_filter := bson.M{"teamId": team_id}
	project_update := bson.M{"$set": bson.M{"isfinal": true}}
	project_result, err := projectCollection.UpdateOne(context.Background(), project_filter, project_update)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "project": project_result})
}

func (databaseClient Database) GetStatus(ctx *fiber.Ctx) error {

	team_id := databaseClient.GetTeamIdByUserId(ctx)
	id_message := *team_id
	if id_message == "ID does not exist" || id_message == "User id not parseable" || id_message == "user not found" || id_message == "user not in a team" {
		fmt.Println("error: ", *team_id)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": *team_id})
	}

	//getting project
	projectCollection := databaseClient.MongoClient.Database("devsoc").Collection("projects")
	findProject := models.Project{}
	project_filter := bson.M{"teamId": team_id}

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

func (databaseClient Database) LikeProject(ctx *fiber.Ctx) error {
	idString := ctx.GetRespHeader("currentUser")

	if idString == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "No Id passed"})
	}

	userId, err := primitive.ObjectIDFromHex(idString)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User ID not parsable"})
	}
	var projectCollection = databaseClient.MongoClient.Database("devsoc").Collection("projects")
	// var projects []models.Project

	projectId := ctx.Params("projectId")

	findProject := models.Project{}
	findTeam := models.Team{}
	prjId, err := primitive.ObjectIDFromHex(projectId)

	projectFilter := bson.M{"_id": prjId}

	err = projectCollection.FindOne(context.TODO(), projectFilter).Decode(&findProject)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Project not found"})
	}
	likedUser := false
	for i := 0; i < findProject.LikeCount; i++ {
		if userId == findProject.LikesId[i] {
			likedUser = true
			fmt.Println("hehe", findProject.LikesId[i])
		}
	}
	if likedUser {
		newDCount := findProject.LikeCount - 1
		projectDUpdate := bson.M{"$set": bson.M{"likecount": newDCount}, "$pull": bson.M{"likesid": userId}}
		projectDRes, err := projectCollection.UpdateOne(context.Background(), projectFilter, projectDUpdate)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "failed to dislike"})
		}
		fmt.Print(projectDRes)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "you already liked:: so -1 dislike"})
	}

	teamCollection := databaseClient.MongoClient.Database("devsoc").Collection("teams")
	teamid := *findProject.TeamId
	teammid, err := primitive.ObjectIDFromHex(teamid)
	teamFilter := bson.M{"_id": teammid}
	err = teamCollection.FindOne(context.TODO(), teamFilter).Decode(&findTeam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "failed to get team"})
	}
	if findTeam.Round < 3 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "err": "Team project is not in 3rd round to vote"})
	}
	newCount := findProject.LikeCount + 1
	projectUpdate := bson.M{"$set": bson.M{"likecount": newCount}, "$push": bson.M{"likesid": userId}}
	projectRes, err := projectCollection.UpdateOne(context.Background(), projectFilter, projectUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "failed to like"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "+1 liked", "doc": projectRes})
}
