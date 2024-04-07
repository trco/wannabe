package actions

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// executes request and sets response fields in ctx *fiber.Ctx to the values received in response
func FetchResponse(ctx *fiber.Ctx, server string) error {
	url := server + ctx.OriginalURL()

	err := proxy.Do(ctx, url)
	if err != nil {
		return fmt.Errorf("FetchResponse: failed executing request: %v", err)
	}

	return nil
}
