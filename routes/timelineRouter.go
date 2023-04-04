package routes

import (
	controller "devsoc23-backend/controllers"
	"devsoc23-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func TimelineRoutes(incomingRoutes *fiber.App, h *controller.Database) {
	incomingRoutes.Get("/timeline", h.GetAllTimeLine)
	incomingRoutes.Get("/timeline/:id", h.GetTimeLine)
	adminGroup := incomingRoutes.Group("/admin", middleware.VerfiyAdmin)

	adminGroup.Post("/timeline/create", h.CreateTimeLine)
	adminGroup.Patch("/timeline/:id", h.UpdateTimeLine)
	adminGroup.Delete("/timeline/:id", h.DeleteTimeline)

}