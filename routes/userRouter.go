package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	incomingRoutes.Post("/users/login", h.LoginUser)
	incomingRoutes.Post("/users/signup", h.RegisterUser)
	incomingRoutes.Get("/users", h.GetUsers)
	incomingRoutes.Post("/users/refresh", h.RefreshToken)
	incomingRoutes.Post("/users/verify", h.Verifyotp)
	incomingRoutes.Post("/users/otp", h.Sendotp)
	incomingRoutes.Post("/users/forgot/mail", h.ForgotPasswordMail)
	incomingRoutes.Post("/users/forgot", h.ForgotPassword)

	userGroup := incomingRoutes.Group("/users", middleware.VerifyToken)
	userGroup.Get("/me", h.FindUser)
	userGroup.Post("/update", h.UpdateUser)
	userGroup.Delete("/delete", h.DeleteUser)
	userGroup.Get("/logout", h.LogoutUser)
	userGroup.Get("/reset", h.ResetPassword)

	adminGroup := incomingRoutes.Group("/admin", middleware.VerfiyAdmin)
	adminGroup.Post("/checkin/:userId", h.CheckIn)
	adminGroup.Post("/checkout/:userId", h.CheckOut)
}
