package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func TeamRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	teamGroup := incomingRoutes.Group("/team", middleware.VerifyToken)
	teamGroup.Post("/create", h.CreateTeam)
	teamGroup.Post("/join/:teamId/:inviteCode", h.JoinTeam)
	teamGroup.Post("/leave", h.LeaveTeam)
	teamGroup.Get("/get/:teamId", h.GetTeam)
	teamGroup.Get("/all", h.GetTeams)
	teamGroup.Get("/members/:teamId", h.GetTeamMembers)
	teamGroup.Get("/ismember", h.GetIsMember)
	teamGroup.Post("/:teamId", h.UpdateTeam)
	teamGroup.Delete("/:teamId", h.DeleteTeam)
	teamGroup.Post("/remove/:teamId/:memberId", h.RemoveMember)

	adminGroup := incomingRoutes.Group("/admin", middleware.VerfiyAdmin)
	adminGroup.Post("/promote/:teamId", h.PromoteTeam)
	adminGroup.Post("/finalise/:teamId", h.FinaliseTeam)
	adminGroup.Post("/disqualify/:teamId", h.DisqualifyTeam)

}
