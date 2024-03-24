package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
)

// set statusCode, headers and body from record to ctx *fiber.Ctx
func SetResponse(ctx *fiber.Ctx, recordBytes []byte) error {
	var record entities.Record

	err := json.Unmarshal(recordBytes, &record)
	if err != nil {
		return fmt.Errorf("SetResponse: failed unmarshaling record: %v", err)
	}

	ctx.Status(record.Response.StatusCode)
	setHeaders(ctx, record.Response.Headers)
	ctx.JSON(record.Response.Body)

	return nil
}

func setHeaders(ctx *fiber.Ctx, headers map[string][]string) {
	for key, value := range headers {
		for _, v := range value {
			ctx.Set(key, v)
		}
	}
}
