package actions

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestGenerateRecord(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	ctx.Request().Header.SetMethod("POST")
	ctx.Method("POST")
	ctx.Path("/test")
	ctx.Request().URI().SetQueryString("test")
	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().Header.Set("Accept", "test")
	// {"test":"test"}
	ctx.Request().SetBody([]byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125})

	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Response().Header.Set("Accept", "test")
	// {"test":"test"}
	ctx.Response().SetBody([]byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125})

	expected, _ := json.Marshal(entities.Record{
		Request: entities.Request{
			HttpMethod: "POST",
			Host:       "test.com",
			Path:       "/test",
			Query: map[string]string{
				"test": "test",
			},
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
			Hash: "123",
			Curl: "test",
		},
		Response: entities.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		},
	})

	record, _ := GenerateRecord(ctx, "test.com", "123", "test")

	if !reflect.DeepEqual(expected, record) {
		t.Errorf("Expected record: %v, Actual record: %v", expected, record)
	}
}
