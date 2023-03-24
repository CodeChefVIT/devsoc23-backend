package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	// incomingRoutes.GET("/users", controller.GetUsers())
	// incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.Post("/users/login", h.LoginUser)
	incomingRoutes.Post("/users/signup", h.RegisterUser)
	incomingRoutes.Get("/users", h.GetUsers)
	incomingRoutes.Post("/users/refresh", h.RefreshToken)
	incomingRoutes.Get("/users/sendotp", h.Sendotp)
	incomingRoutes.Patch("/users/verifyotp", h.Verifyotp)

	userGroup := incomingRoutes.Group("/users", middleware.VerifyToken)
	userGroup.Get("/me", h.FindUser)
	userGroup.Get("/logout", h.LogoutUser)
	userGroup.Get("/reset", h.ResetPassword)

	userGroup.Patch("/checkin", h.CheckIn)
	userGroup.Patch("/checkout", h.CheckOut)

}
