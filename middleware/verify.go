package middleware

import (
	"context"
	"devsoc23-backend/database"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyToken(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	filter := bson.M{"_id": res.Id}

	user := models.User{}
	client := database.NewDatabase().MongoClient
	userCollection := client.Database("devsoc").Collection("users")

	err = userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists", "user": user, "err": err.Error(), "jwt": res})
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// TODO optional you can log your closed MongoDB client
	// fmt.Println("Connection to MongoDB closed.")

	ctx.Set("currentUser", user.Id.Hex())
	return ctx.Next()
}

func VerfiyAdmin(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "ADMIN" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an admin"})
	}

	filter := bson.M{"_id": res.Id}

	user := models.User{}
	client := database.NewDatabase().MongoClient
	userCollection := client.Database("devsoc").Collection("users")

	err = userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	ctx.Set("currentUser", user.Id.Hex())
	return ctx.Next()
}

func VerifyBoard(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	res, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if res.Role != "ADMIN" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not an admin"})
	}

	filter := bson.M{"_id": res.Id}

	user := models.User{}
	client := database.NewDatabase().MongoClient
	userCollection := client.Database("devsoc").Collection("users")

	err = userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
	}

	if !user.IsBoard {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User is not a board member"})
	}

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	ctx.Set("currentUser", user.Id.Hex())
	return ctx.Next()
}
