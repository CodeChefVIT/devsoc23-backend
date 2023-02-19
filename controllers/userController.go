package controller

import (
	"fmt"
	"log"
	"strings"

	"devsoc23-backend/helper"
	models "devsoc23-backend/models"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

// 		}
// 		c.JSON(http.StatusOK, user)
// 	}
// }

type registrationUserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
type registrationResponse struct {
	Token string `json:"token"`
}

// type UserSession struct {
// 	UserID        int
// 	Authenticated bool
// }

func (databaseClient Database) Register(ctx *fiber.Ctx) error {

	var input registrationUserRequest
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	r := databaseClient.DB.Where("email_id = ?", input.Email).Limit(1).Find(&models.User{})
	exists := r.RowsAffected > 0
	if exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User Already Exists in db",
		})

	}

	input.Email = strings.Trim(input.Email, " ")
	if input.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email empty",
		})

	}

	input.Password = strings.Trim(input.Password, " ")
	if input.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password Empty",
		})

	}

	input.FirstName = strings.Trim(input.FirstName, " ")
	if input.FirstName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "First Name empty",
		})

	}

	input.LastName = strings.Trim(input.LastName, " ")
	if input.LastName == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Last Name empty",
		})

	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password error",
		})

	}

	user := models.User{EmailID: input.Email, Password: string(hash), FirstName: input.FirstName, LastName: input.LastName}
	resultUser := databaseClient.DB.Create(&user)
	if resultUser.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something Went Wrong User not created",
		})

	}
	sessionID := helper.GenerateToken()

	err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), sessionID, user.ID, redis.KeepTTL).Err()
	if err != nil {
		databaseClient.DB.Unscoped().Delete(&user)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})

	}
	err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), strconv.FormatUint(uint64(user.ID), 10), sessionID, redis.KeepTTL).Err()
	if err != nil {
		databaseClient.DB.Unscoped().Delete(&user)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to add user session to redis",
		})

	}
	claims := jwt.MapClaims{
		"sessionID": sessionID,
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)

	}
	return ctx.Status(fiber.StatusOK).JSON(registrationResponse{
		Token: t,
	})

}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type loginResponse struct {
	Token string `json:"token"`
}

func (databaseClient Database) Login(ctx *fiber.Ctx) error {
	var input loginUserRequest
	if err := ctx.BodyParser(&input); err != nil {

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var user models.User = models.User{EmailID: input.Email}
	r := databaseClient.DB.Where("email_id = ?", input.Email).Limit(1).Find(&user)
	exists := r.RowsAffected == 0
	if exists {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User does not exist",
		})

	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password did not matched",
		})

	}
	sessionID, err := databaseClient.RedisClient.Get(databaseClient.RedisClient.Context(), strconv.FormatUint(uint64(user.ID), 10)).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UserID does not exist",
		})

	}
	// logout condition
	if sessionID == "" {
		sessionID = helper.GenerateToken()
		err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), sessionID, user.ID, redis.KeepTTL).Err()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "unable to add user session to redis",
			})

		}
		err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), strconv.FormatUint(uint64(user.ID), 10), sessionID, redis.KeepTTL).Err()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "unable to add user session to redis",
			})

		}
	}
	claims := jwt.MapClaims{
		"sessionID": sessionID,
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)

	}

	return ctx.Status(fiber.StatusOK).JSON(loginResponse{Token: t})

}
func (databaseClient Database) Logout(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	sessionToken := claims["sessionID"].(string)

	userID, err := databaseClient.RedisClient.Get(databaseClient.RedisClient.Context(), sessionToken).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": sessionToken,
		})

	}
	_, err = databaseClient.RedisClient.Del(databaseClient.RedisClient.Context(), sessionToken).Result()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to remove user token from redis",
		})

	}
	err = databaseClient.RedisClient.Set(databaseClient.RedisClient.Context(), userID, "", redis.KeepTTL).Err()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to add user token to redis",
		})

	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logged out",
	})

}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false

	}
	return check, msg
}
