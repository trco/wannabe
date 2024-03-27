package handlers

import "github.com/gofiber/fiber/v2"

func internalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(InternalError{
		StatusCode: fiber.StatusInternalServerError,
		Error:      err.Error(),
	})
}

func checkDuplicates(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

func processRecordValidation(recordProcessingDetails *[]RecordProcessingDetails, hash string, message string, valueToIncrement *int) {
	*recordProcessingDetails = append(*recordProcessingDetails, RecordProcessingDetails{
		Hash:    hash,
		Message: message,
	})

	*valueToIncrement++
}
