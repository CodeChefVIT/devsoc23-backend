package middleware

import (
	"context"
	"devsoc23-backend/database"
	"devsoc23-backend/helper"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	helper.LoadEnv()
}

func VerifyToken(ctx *fiber.Ctx) error {

	var token string

	authorizationHeader := ctx.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	sub, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	filter := bson.D{{"email", sub}}

	user := models.NewUser{}
	userCollection := database.NewDatabase().MongoClient.Database("devsoc").Collection("users")

	err = userCollection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no longer exists"})
	}

	ctx.Set("currentUser", *user.Email)
	return ctx.Next()

}
