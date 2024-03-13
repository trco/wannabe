package handlers

import "github.com/gofiber/fiber/v2"

func Readiness(ctx *fiber.Ctx) error {
	return ctx.SendString("I'm ready!")
}
