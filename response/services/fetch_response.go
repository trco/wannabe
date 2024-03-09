package services

import (
	"wannabe/response/actions"

	"github.com/gofiber/fiber/v2"
)

func FetchResponse(ctx *fiber.Ctx, server string) error {
	err := actions.FetchResponse(ctx, server)
	if err != nil {
		return err
	}

	return nil
}
