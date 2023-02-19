package routes

import (
	controller "devsoc23-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	// incomingRoutes.GET("/users", controller.GetUsers())
	// incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.Post("/users/signup", h.Register)
	// incomingRoutes.POST("/users/login", controller.Login())
}
