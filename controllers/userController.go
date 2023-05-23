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

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var store = make(map[string]string)
var reset = make(map[string]string)

type EmailData struct {
	Email string `json:"email"`
}

type EmailOTPData struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
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
	var payload *models.CreateUserRequest

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	filter := bson.M{"email": payload.Email}
	count, _ := userCollection.CountDocuments(context.TODO(), filter)
	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "Email already exists"})
	}

	now := time.Now()
	userRole := "HACKER"
	hash, _ := utils.HashPassword(*payload.Password)

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Password:    &hash,
		PhoneNumber: payload.PhoneNumber,
		UserRole:    userRole,
		Bio:         payload.Bio,
		Gender:      payload.Gender,
		RegNo:       payload.RegNo,
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
		Id:   newUser.Id,
		Role: newUser.UserRole,
	}
	token, err := utils.GenerateToken(duration, sub, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	// Update refreshToken in user document
	err = databaseClient.RedisClient.Set(context.Background(), token, newUser.Id.String(), 0).Err()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": result})
}

func (databaseClient Database) FindUser(ctx *fiber.Ctx) error {

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	id, err := primitive.ObjectIDFromHex(ctx.GetRespHeader("currentUser"))

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email does not exist"})
	}

	findUser := models.User{}
	filter := bson.M{"_id": id}

	err = userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}
	fmt.Println(findUser.Id)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": findUser})
}

func (databaseClient Database) UpdateUser(ctx *fiber.Ctx) error {

	// Get request body and bind to payload
	var payload models.UpdateUserRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	file, err := ctx.FormFile("image")
	url := ""
	if err == nil {
		image := utils.PhotoForm{
			CampaignImage: file,
		}
		url, uploadErr := utils.UploadPhoto(&image, databaseClient.S3Client)
		fmt.Println(url)
		if uploadErr != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err.Error(), "message": "Image Upload Failed"})
		}
	}

	// Get current user from the response header
	user := ctx.GetRespHeader("currentUser")
	id, err := primitive.ObjectIDFromHex(user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	// Create user object
	update := bson.M{
		"firstname":   payload.FirstName,
		"lastname":    payload.LastName,
		"email":       payload.Email,
		"phonenumber": payload.PhoneNumber,
		"college":     payload.College,
		"bio":         payload.Bio,
		"gender":      payload.Gender,
		"regno":       payload.RegNo,
		"collegeyear": payload.CollegeYear,
		"birthdate":   payload.BirthDate,
		"image":       url,
		"isactive":    false,
		"iscanshare":  false,
		"updatedat":   time.Now(),
	}

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	// Update user in user document
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "User updated succesfully"})
}

func (databaseClient Database) DeleteUser(ctx *fiber.Ctx) error {

	// Get current user from the response header
	user := ctx.GetRespHeader("currentUser")
	id, err := primitive.ObjectIDFromHex(user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")

	// Delete user from user document
	filter := bson.M{"_id": id}
	deleteResult, _ := userCollection.DeleteOne(context.TODO(), filter)

	if deleteResult.DeletedCount == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not deleted"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "User deleted succesfully"})
}

func (databaseClient Database) RefreshToken(ctx *fiber.Ctx) error {

	type tokenRequest struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	payload := tokenRequest{}
	// Get refreshToken from request
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Find refresh token in redis
	_, err := databaseClient.RedisClient.Get(context.Background(), payload.RefreshToken).Result()
	if err != nil {
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
	match := utils.CheckPasswordHash(*payload.Password, *findUser.Password)

	if !match {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Wrong password"})
	}

	// Create a new refreshToken
	duration, _ := time.ParseDuration("1h")
	sub := utils.TokenPayload{
		Id:   findUser.Id,
		Role: findUser.UserRole,
	}
	token, err := utils.GenerateToken(duration, sub, os.Getenv("REFRESH_JWT_SECRET"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": err})
	}

	// Update refreshToken in user document
	err = databaseClient.RedisClient.Set(context.Background(), token, findUser.Id.String(), 0).Err()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update refreshToken"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "user": findUser, "token": token})
}

func (databaseClient Database) LogoutUser(ctx *fiber.Ctx) error {

	type tokenRequest struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	payload := tokenRequest{}
	// Get refreshToken from request
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Delete refresh token in redis
	_, err := databaseClient.RedisClient.Del(context.Background(), payload.RefreshToken).Result()

	fmt.Println(err)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not logout user"})
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
	subject := "Your OTP for Devsoc verification"
	body := "Your Devsoc verification OTP is: " + otp
	err := utils.SendMail(subject, body, email)

	// Send the email
	if err != nil {
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

func (databaseClient Database) ResetPassword(ctx *fiber.Ctx) error {

	// Get request body and bind to payload
	var payload *models.ResetPasswordRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}

	// Validate Struct
	errors := utils.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Check if oldpass and newpass are same

	hash, _ := utils.HashPassword(*payload.Newpass)
	match := utils.CheckPasswordHash(*payload.Oldpass, hash)

	if match {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": "Old password and New password cannot be the same."})
	}

	// Get User
	userCollection := databaseClient.MongoClient.Database("devsoc").Collection("users")
	user := ctx.GetRespHeader("currentUser")
	id, err := primitive.ObjectIDFromHex(user)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User does not exist"})
	}

	findUser := models.User{}
	filter := bson.M{"_id": id}

	err = userCollection.FindOne(context.TODO(), filter).Decode(&findUser)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "User not found"})
	}

	// Check if oldpass matches

	match = utils.CheckPasswordHash(*payload.Oldpass, *findUser.Password)

	if !match {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Old password is incorrect"})
	}

	// Update user with new password
	_, err = userCollection.UpdateOne(context.TODO(), bson.M{"email": findUser.Email}, bson.M{"$set": bson.M{"password": &hash}})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Could not update password"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Password successfully reset"})
}

func (databaseClient Database) ForgotPasswordMail(ctx *fiber.Ctx) error {
	var userCollection = databaseClient.MongoClient.Database("devsoc").Collection("users")
	// Get request body and bind to payload
	var payload EmailData
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}
	email := payload.Email
	findUser := models.User{}
	filter := bson.M{"email": email}
	errr := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email not found"})
	}

	// Generate OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store OTP
	reset[email] = otp

	// Send email with OTP
	subject := "Password Reset Request"
	body := "Your password reset OTP is: " + otp
	err := utils.SendMail(subject, body, email)

	// Send the email
	if err != nil {
		// Return error if email cannot be sent
		fmt.Println("Error reason: ", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": "Failed to send OTP",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Reset email sent successfully"})
}

func (databaseClient Database) ForgotPassword(ctx *fiber.Ctx) error {
	var userCollection = databaseClient.MongoClient.Database("devsoc").Collection("users")
	// Get request body and bind to payload
	var payload *models.ForgetPasswordRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "err": err.Error()})
	}
	email := payload.Email
	otp := payload.OTP
	findUser := models.User{}
	filter := bson.M{"email": email}
	errr := userCollection.FindOne(context.TODO(), filter).Decode(&findUser)
	if errr != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "false", "err": "Email not found"})
	}
	// Retrieve OTP from store
	storedOtp, ok := reset[email]
	if !ok {
		// Return error if OTP not found for email
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "OTP not found",
		})
	}

	// Compare OTP
	if otp != storedOtp {
		// Return error if OTP is incorrect
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Incorrect OTP",
		})
	}

	// Delete OTP from store
	delete(store, email)

	hash, _ := utils.HashPassword(payload.Newpass)
	// Create a filter to find the user with the given email
	// Update new password
	update := bson.M{"$set": bson.M{"password": hash}}

	// Update the user record in the database
	result, err := userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// User not found
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "user not found"})

	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true", "message": "Password updated successfully"})
}
