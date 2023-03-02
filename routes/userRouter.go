package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, h *controller.Database) {
	
	// incomingRoutes.GET("/users", controller.GetUsers())
	// incomingRoutes.GET("/users/:user_id", controller.GetUser())
	// incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.Post("/users/signup", h.RegisterUser)
	protected := incomingRoutes.Group("/users", middleware.VerifyToken)
	protected.Get("/me",h.FindUser)
}
