package actions

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// executes request and sets response fields in ctx *fiber.Ctx to the values received in response
func FetchResponse(ctx *fiber.Ctx, protocol string, host string) error {
	fmt.Println("ctx.OriginalURL", ctx.OriginalURL())

	uri := string(ctx.Request().RequestURI())

	fmt.Println("uri", uri)

	url := protocol + "://" + host + uri
	// url := protocol + "://" + host + ctx.OriginalURL()

	fmt.Println("url", url)

	err := proxy.Do(ctx, url)
	if err != nil {
		return fmt.Errorf("FetchResponse: failed executing request: %v", err)
	}

	return nil
}
