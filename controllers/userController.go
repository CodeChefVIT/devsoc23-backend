package controller

import (
	// "context"
	// "fmt"

	"context"
	"fmt"
	"os"

	// "strconv"

	helper "devsoc23-backend/helper"
	models "devsoc23-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"

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

// 	}
// }
// func GetUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		userId := c.Param("user_id")
// 		var user models.User
// 		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user"})

//			}
//			c.JSON(http.StatusOK, user)
//		}
//	}

var validate = validator.New()

type registrationUserRequest struct {
	FirstName   *string `json:"firstName" validate:"required,min=2,max=16"`
	LastName    *string `json:"lastName" validate:"min=2,max=32"`
	Email       *string `json:"email" validate:"required,email"`
	Password    *string `json:"password" validate:"required,min=8,max=64"`
	PhoneNumber *string `json:"phoneNumber" validate:"required,e164,min=10,max=13"`
}
type registrationResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(user registrationUserRequest) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
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

	var input registrationUserRequest
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

	user := models.User{Email: input.Email, Password: &hash, FirstName: input.FirstName, LastName: input.LastName, PhoneNumber: input.PhoneNumber}
	registrationResponse, errr := userCollection.InsertOne(context.Background() , user)
	if errr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something Went Wrong User not created",
		})
	}
	fmt.Print(registrationResponse)

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
		"token": t,
	})

}

// type loginUserRequest struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }
// type loginResponse struct {
// 	Token string `json:"token"`
// }

// func (databaseClient Database) Login(ctx *fiber.Ctx) error {
// 	var input loginUserRequest
// 	if err := ctx.BodyParser(&input); err != nil {

// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	var user models.User = models.User{EmailID: input.Email}
// 	r := databaseClient.DB.Where("email_id = ?", input.Email).Limit(1).Find(&user)
// 	exists := r.RowsAffected == 0
// 	if exists {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "User does not exist",
// 		})

// 	}

// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Password did not matched",
// 		})

// 	}
// 	sessionID, err := databaseClient.RedisClient.Get(databaseClient.RedisClient.Context(), strconv.FormatUint(uint64(user.ID), 10)).Result()
// 	if err == redis.Nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "UserID does not exist",
// 		})

// 	}
// 	// logout condition
// 	if sessionID == "" {
// 		sessionID = helper.GenerateToken()
// 		err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), sessionID, user.ID, redis.KeepTTL).Err()
// 		if err != nil {
// 			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "unable to add user session to redis",
// 			})

// 		}
// 		err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), strconv.FormatUint(uint64(user.ID), 10), sessionID, redis.KeepTTL).Err()
// 		if err != nil {
// 			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "unable to add user session to redis",
// 			})

// 		}
// 	}
// 	claims := jwt.MapClaims{
// 		"sessionID": sessionID,
// 	}
// 	// Create token
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Generate encoded token and send it as response.
// 	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	if err != nil {
// 		ctx.SendStatus(fiber.StatusInternalServerError)

// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(loginResponse{Token: t})

// }
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

// func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
// 	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
// 	check := true
// 	msg := ""
// 	if err != nil {
// 		msg = fmt.Sprintf("login or password is incorrect")
// 		check = false

// 	}
// 	return check, msg
// }
