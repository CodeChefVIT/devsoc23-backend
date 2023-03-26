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
	incomingRoutes.Patch("/users/verifyotp", h.Verifyotp)
	incomingRoutes.Post("/sendotp", h.Sendotp)
	incomingRoutes.Post("/users/forgot/mail", h.ForgotPasswordMail)
	incomingRoutes.Patch("/users/forgot", h.ForgotPassword)

	
	userGroup := incomingRoutes.Group("/users", middleware.VerifyToken)
	userGroup.Get("/me", h.FindUser)
	userGroup.Patch("/update", h.UpdateUser)
	userGroup.Delete("/delete", h.DeleteUser)
	userGroup.Get("/logout", h.LogoutUser)
	userGroup.Get("/reset", h.ResetPassword)
	
	adminGroup := incomingRoutes.Group("/admin", middleware.VerfiyAdmin)
	adminGroup.Get("/me", h.FindUser)
	adminGroup.Patch("/checkin", h.CheckIn)
	adminGroup.Patch("/checkout", h.CheckOut)

}
