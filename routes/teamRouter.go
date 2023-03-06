package routes

import (
	controller "devsoc23-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func TeamRoutes(incomingRoutes *fiber.App, h *controller.Database) {
	/*
		teamGroup := incomingRoutes.Group("/teams", middleware.VerifyToken)
		teamGroup.Post("/create", h.CreateTeam)
		teamGroup.Get("/:teamId", h.GetTeam)
		teamGroup.Get("/all", h.GetTeams)
		teamGroup.Post("/:teamId", h.UpdateTeam)
		teamGroup.Delete("/:teamId", h.DeleteTeam)

		teamGroup.Patch("/:teamId", h.JoinTeam)
		teamGroup.Patch("/:teamId", h.FinaliseTeam)
	*/
}
