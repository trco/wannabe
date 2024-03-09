package handlers

import "github.com/gofiber/fiber/v2"

func Liveness(c *fiber.Ctx) error {
	return c.SendString("I'm alive!")
}
