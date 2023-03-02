package controller

import (
	// "context"
	// "fmt"

	"context"
	// "fmt"
	"os"
	"time"

	// "strconv"

	helper "devsoc23-backend/helper"
	models "devsoc23-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"

	// "golang.org/x/crypto/bcrypt"
	// "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	// "go.mongodb.org/mongo-driver/mongo"
)

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}
// 		page, err1 := strconv.Atoi(c.Query("page"))
// 		if err1 != nil || page < 1 {
// 			page = 10
// 		}
// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err = strconv.Atoi(c.Query("startIndex"))
// 		matchStage := bson.D{{"$match", bson.D{{}}}}
// 		projectStage := bson.D{
// 			{"$project", bson.D{
// 				{"_id", 0},
// 				{"total_count", 1},
// 				{"user_item", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
// 			}},
// 		}
// 		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage,
// 			projectStage,
// 		})
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
// 		}
// 		var allUsers []bson.M
// 		if err = result.All(ctx, &allUsers); err != nil {
// 			log.Fatal(err)

// 		}
// 		c.JSON(http.StatusOK, allUsers)

//		}
//	}
func (databaseClient Database) GetUser(ctx *fiber.Ctx) error {

	// userCollection := databaseClient.MongoClient.Database("devsoc23").Collection("user")
	var userCollection *mongo.Collection = databaseClient.MongoClient.Database("devsoc23").Collection("user")
	userId := ctx.Query("userId")
	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"userid": userId}).Decode(&user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}

var validate = validator.New()

func ValidateStruct(user models.RegistrationUserRequest) []*models.ErrorResponse {
	var errors []*models.ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element models.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (databaseClient Database) Singup(ctx *fiber.Ctx) error {

	userCollection := databaseClient.MongoClient.Database("devsoc23").Collection("user")

	var input models.RegistrationUserRequest
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := ValidateStruct(input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	count, errr := userCollection.CountDocuments(context.Background(), bson.M{"email": input.Email})
	if errr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errr.Error(),
		})
	}

	if count > 0 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "email is already taken",
		})
	}

	hash := helper.HashPassword(*input.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something Went Wrong User not created",
		})
	}

	role := "HACKER"

	user := models.User{Email: input.Email, Password: &hash, FirstName: input.FirstName, LastName: input.LastName, PhoneNumber: input.PhoneNumber, IsActive: true, IsVerify: false, IsBoard: false, IsCanShare: false, IsCheckedIn: false, InTeam: false}
	user.Id = primitive.NewObjectID()
	user.UserId = user.Id.Hex()
	user.UserRole = &role
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	RegistrationResponse, errr := userCollection.InsertOne(context.Background(), user)
	if errr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something Went Wrong User not created",
		})
	}
	sessionID := helper.GenerateToken()

	claims := jwt.MapClaims{
		"sessionID": sessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, errr := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errr != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    RegistrationResponse.InsertedID,
		"token": t,
	})

}

func (databaseClient Database) Login(ctx *fiber.Ctx) error {
	var user models.User
	var input models.LoginUserRequest

	userCollection := databaseClient.MongoClient.Database("devsoc23").Collection("user")

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := userCollection.FindOne(context.Background(), bson.M{"email": *input.Email}).Decode(&input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	pssderr := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*input.Password))
	if pssderr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password did not matched",
		})

	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"sucess": "logged in succesffuly",
	})
}

// func (databaseClient Database) Logout(ctx *fiber.Ctx) error {
// 	user := ctx.Locals("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	sessionToken := claims["sessionID"].(string)

// 	userID, err := databaseClient.RedisClient.Get(databaseClient.RedisClient.Context(), sessionToken).Result()
// 	if err == redis.Nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": sessionToken,
// 		})

// 	}
// 	_, err = databaseClient.RedisClient.Del(databaseClient.RedisClient.Context(), sessionToken).Result()
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "unable to remove user token from redis",
// 		})

// 	}
// 	err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), userID, "", redis.KeepTTL).Err()
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "unable to add user token to redis",
// 		})

// 	}
// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "logged out",
// 	})

// }
