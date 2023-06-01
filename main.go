package main

import (
	"devsoc23-backend/database"
	helper "devsoc23-backend/helper"
	"devsoc23-backend/routes"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	helper.LoadEnv()
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Api rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        500,
		Expiration: 30 * time.Second,
	}))
	app.Use(logger.New())

	handler := database.NewDatabase()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})
	routes.UserRoutes(app, &handler)
	routes.TeamRoutes(app, &handler)
	routes.ProjectRoutes(app, &handler)
	routes.TimelineRoutes(app, &handler)
	err = app.Listen(":8000")
	if err != nil {
		panic(err)
	}
}
