package services

import (
	"wannabe/response/actions"

	"github.com/gofiber/fiber/v2"
)

func SetResponse(ctx *fiber.Ctx, recordBytes []byte) error {
	err := actions.SetResponse(ctx, recordBytes)
	if err != nil {
		return err
	}

	return nil
}
