package handlers

import "github.com/gofiber/fiber/v2"

type WannabeHandler func(ctx *fiber.Ctx) error

type InternalError struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

func internalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(InternalError{
		StatusCode: fiber.StatusInternalServerError,
		Error:      err.Error(),
	})
}

type PostRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type DeleteRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type RegenerateResponse struct {
	Message           string   `json:"message"`
	RegeneratedHashes []string `json:"regeneratedHashes"`
	FailedHashes      []string `json:"failedHashes"`
}
