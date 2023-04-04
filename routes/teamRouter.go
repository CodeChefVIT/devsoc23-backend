package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func TeamRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	teamGroup := incomingRoutes.Group("/team", middleware.VerifyToken)
	teamGroup.Post("/create", h.CreateTeam)
	teamGroup.Patch("/join/:teamId/:inviteCode", h.JoinTeam)
	teamGroup.Patch("/leave", h.LeaveTeam)
	teamGroup.Get("/get/:teamId", h.GetTeam)
	teamGroup.Get("/all", h.GetTeams)
	teamGroup.Get("/members/:teamId", h.GetTeamMembers)
	teamGroup.Post("/:teamId", h.UpdateTeam)
	teamGroup.Delete("/:teamId", h.DeleteTeam)
	
	adminGroup := incomingRoutes.Group("/admin", middleware.VerfiyAdmin)
	adminGroup.Patch("/promote/:teamId", h.PromoteTeam)
	adminGroup.Patch("/finalise/:teamId", h.FinaliseTeam)
	adminGroup.Patch("/disqualify/:teamId",h.DisqualifyTeam)

}
