package services

import (
	"wannabe/response/actions"

	"github.com/gofiber/fiber/v2"
)

func FetchResponse(ctx *fiber.Ctx, protocol string, host string) error {
	err := actions.FetchResponse(ctx, protocol, host)
	if err != nil {
		return err
	}

	return nil
}
