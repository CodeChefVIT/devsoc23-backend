package controller

import (
	"github.com/gofiber/fiber/v2"
)

func (databaseClient Database) CreateTimeLine(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true"})
}

func (databaseClient Database) GetAllTimeLine(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true"})
}
func (databaseClient Database) GetTimeLine(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true"})
}
func (databaseClient Database) UpdateTimeLine(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "true"})
}
