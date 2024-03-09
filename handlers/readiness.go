package handlers

import "github.com/gofiber/fiber/v2"

func Readiness(c *fiber.Ctx) error {
	return c.SendString("I'm ready!")
}
