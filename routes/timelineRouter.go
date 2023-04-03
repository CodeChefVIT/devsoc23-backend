package routes

import (
	controller "devsoc23-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func TimelineRoutes(incomingRoutes *fiber.App, h *controller.Database) {

	incomingRoutes.Post("/timeline/login", h.CreateTimeLine)
	incomingRoutes.Get("/timeline", h.GetAllTimeLine)
	incomingRoutes.Get("/timeline/:id", h.GetTimeLine)
	incomingRoutes.Patch("/timeline/:id", h.UpdateTimeLine)

}
