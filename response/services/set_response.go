package services

import (
	"wannabe/response/actions"

	"github.com/gofiber/fiber/v2"
)

func SetResponse(ctx *fiber.Ctx, encodedRecord []byte) error {
	err := actions.SetResponse(ctx, encodedRecord)
	if err != nil {
		return err
	}

	return nil
}
