package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func ProjectRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	projectGroup := incomingRoutes.Group("/project", middleware.VerifyToken)
	projectGroup.Post("/idea",h.CreateProjectIdea)
	// projectGroup.Post("/create", h.CreateProject)
	projectGroup.Get("/get", h.GetProjectByUserid)
	projectGroup.Get("/get/:teamId", h.GetProjectByTeamid)
	projectGroup.Post("/update", h.UpdateProject)
	projectGroup.Delete("/delete", h.DeleteProject)
	projectGroup.Post("/finalproject", h.FinaliseProjectSubmission)
	projectGroup.Get("/status", h.GetStatus)
	projectGroup.Get("/allprojects", h.GetProjects)

	projectGroup.Patch("/like/:projectId", h.LikeProject)
}
