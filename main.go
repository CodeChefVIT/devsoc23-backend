package main

import (
	"devsoc23-backend/database"
	helper "devsoc23-backend/helper"
	"devsoc23-backend/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	// Api rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
	}))
	app.Use(logger.New())
	helper.LoadEnv()
	handler := database.NewDatabase()
	//s3Client := infrastructure.InitializeSpaces()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})
	routes.UserRoutes(app, &handler)
	routes.TeamRoutes(app, &handler)
	routes.ProjectRoutes(app, &handler)
	err := app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
