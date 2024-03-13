package handlers

import "github.com/gofiber/fiber/v2"

func Liveness(ctx *fiber.Ctx) error {
	return ctx.SendString("I'm alive!")
}
