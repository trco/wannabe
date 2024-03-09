package services

import (
	"wannabe/record/actions"

	"github.com/gofiber/fiber/v2"
)

// FIXME remove ctx as parameter, extract values and pass them in custom struct
func GenerateRecord(ctx *fiber.Ctx, server string, curl string, hash string) ([]byte, error) {
	record, err := actions.GenerateRecord(ctx, server, curl, hash)
	if err != nil {
		return nil, err
	}

	return record, nil
}
