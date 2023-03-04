package routes

import (
	controller "devsoc23-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	// incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.Get("/users/:userId", h.GetUser)
	incomingRoutes.Post("/auth/signup", h.Singup)
	incomingRoutes.Post("/auth/login", h.Login)
}
