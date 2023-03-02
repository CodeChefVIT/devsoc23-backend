package controller

import (
	"context"
	"devsoc23-backend/helper"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

// type registrationUserRequest struct {
// 	Email     string `json:"email" binding:"required"`
// 	Password  string `json:"password" binding:"required"`
// 	FirstName string `json:"first_name" binding:"required"`
// 	LastName  string `json:"last_name" binding:"required"`
// }
// type registrationResponse struct {
// 	Token string `json:"token"`
// }

//	type UserSession struct {
//		UserID        int
//		Authenticated bool
//	}

func (databaseClient Database) GetUsers(ctx *fiber.Ctx) error {
	var userCollection = databaseClient.MongoClient.Database("devsoc").Collection("users")
	var users []models.User
	cur, err := userCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()

		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)

	}
	if err := cur.Err(); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "users": users})
}

func (databaseClient Database) RegisterUser(ctx *fiber.Ctx) error {
	//Incomplete logic
	var payload *models.CreateUserRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)

	}

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	now := time.Now()
	userRole := "HACKER"
	hash, _ := HashPassword(*payload.Password)

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Password:    &hash,
		PhoneNumber: payload.PhoneNumber,
		UserRole:    userRole,
		College:     payload.College,
		CollegeYear: payload.CollegeYear,
		BirthDate:   payload.BirthDate,
		IsActive:    false,
		IsVerify:    false,
		IsCanShare:  false,
		IsCheckedIn: false,
		InTeam:      false,
		IsBoard:     false,
		CreatedAt:   now,
		UpdatedAt:   now,
		UserId:      helper.GenerateToken(),
	}

	result, err := userCollection.InsertOne(context.TODO(), newUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}
	fmt.Println(result.InsertedID)

	duration, _ := time.ParseDuration("1h")
	token, err := utils.GenerateToken(duration, newUser.Email, os.Getenv("JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "id": result.InsertedID, "token": token})
}

func (databaseClient Database) FindUser(ctx *fiber.Ctx) error {

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	email := ctx.GetRespHeader("currentUser")

	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email does not exist"})
	}

	findUser := models.User{}
	filter := bson.D{{"email", email}}

	err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}
	fmt.Println(findUser.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": findUser})
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

// func HashPassword(password string) string {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	return string(bytes)
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
