package controller

import (
	"context"
	"devsoc23-backend/helper"
	"devsoc23-backend/models"
	"devsoc23-backend/utils"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var store = make(map[string]string)

type EmailData struct {
	Email string `json:"email"`
}

type EmailOTPData struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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
	//Incomplete logic : check for unique user
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
	sub := utils.TokenPayload{
		Email: *newUser.Email,
		Role:  newUser.UserRole,
	}
	token, err := utils.GenerateToken(duration, sub, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	// Update refreshToken in user document
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": newUser.Email}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": result})
}

func (databaseClient Database) FindUser(ctx *fiber.Ctx) error {

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	email := ctx.GetRespHeader("currentUser")

	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email does not exist"})
	}

	findUser := models.User{}
	filter := bson.M{"email": email}

	err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}
	fmt.Println(findUser.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": findUser})
}

func (databaseClient Database) RefreshToken(ctx *fiber.Ctx) error {

	type tokenRequest struct {
		RefreshToken string `json:"refreshToken"`
	}

	payload := tokenRequest{}
	// Get refreshToken from request
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Find refresh token in db
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	filter := bson.M{"token": payload.RefreshToken}

	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if count == 0 || err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "refreshToken not found"})
	}

	// Validate Refresh Token
	sub, err := utils.ValidateToken(payload.RefreshToken, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "err": err.Error(), "token": payload.RefreshToken})
	}

	// Create new accessToken
	duration, _ := time.ParseDuration("1h")
	accessToken, err := utils.GenerateToken(duration, sub, os.Getenv("JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "accessToken": accessToken})
}

func (databaseClient Database) LoginUser(ctx *fiber.Ctx) error {

	// Get request body and bind to payload
	var payload *models.LoginUserRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Find user in collection
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	findUser := models.User{}
	filter := bson.M{"email": payload.Email}

	err := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	// Compare password hashes
	match := CheckPasswordHash(*payload.Password, *findUser.Password)

	if !match {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong password"})
	}

	// Create a new refreshToken
	duration, _ := time.ParseDuration("1h")
	sub := utils.TokenPayload{
		Email: *findUser.Email,
		Role:  findUser.UserRole,
	}
	token, err := utils.GenerateToken(duration, sub, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	// Update refreshToken in user document
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": findUser.Email}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": findUser, "token": token})
}

func (databaseClient Database) LogoutUser(ctx *fiber.Ctx) error {

	// Get current user from the response header
	email := ctx.GetRespHeader("currentUser")

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	// Update refreshToken in user document
	_, err := userCollection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.M{"$set": bson.M{"token": nil}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "User logged out succesfully"})
}

func (databaseClient Database) Sendotp(c *fiber.Ctx) error {
	// Get email from request body
	var emailData EmailData
	if err := c.BodyParser(&emailData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	email := emailData.Email
	// Generate OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store OTP
	store[email] = otp

	// Send email with OTP
	m := gomail.NewMessage()
	m.SetHeader("From", "noreplydevsoc23test@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP for Devsoc verification")
	m.SetBody("text/plain", "Your Devsoc verification OTP is: "+otp)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	from_emailid := os.Getenv("TEST_EMAILID")
	app_pass := os.Getenv(" TEST_EMAILID_PASS")
	d := gomail.NewDialer("smtp.gmail.com", 587, from_emailid, app_pass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		// Return error if email cannot be sent
		fmt.Println("Error reason: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send OTP",
		})
	}
	// Return success message
	return c.JSON(fiber.Map{
		"message": "OTP sent successfully",
	})
}

func (databaseClient Database) Verifyotp(c *fiber.Ctx) error {
	// Get email and OTP from request body
	var emailotpData EmailOTPData
	if err := c.BodyParser(&emailotpData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	email := emailotpData.Email
	otp := emailotpData.OTP
	// Retrieve OTP from store
	storedOtp, ok := store[email]
	if !ok {
		// Return error if OTP not found for email
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "OTP not found",
		})
	}

	// Compare OTP
	if otp != storedOtp {
		// Return error if OTP is incorrect
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Incorrect OTP",
		})
	}

	// Delete OTP from store
	delete(store, email)

	// Return success message
	collection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	// Create a filter to find the user with the given email
	filter := bson.M{"email": email}

	// Create an update document with the new value for the isVerify field
	update := bson.M{"$set": bson.M{"isVerify": true}}

	// Update the user record in the database
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// User not found
		return fmt.Errorf("user with email %s not found", email)
	}
	return c.JSON(fiber.Map{
		"message": "OTP verified successfully",
	})
}
