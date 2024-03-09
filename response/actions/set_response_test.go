package actions

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestSetResponse(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	record, _ := json.Marshal(entities.Record{
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

	expectedStatusCode := 200
	expectedResponseHeaders := map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	}
	expectedResponseBody := map[string]interface{}{
		"test": "test",
	}

	// Invalid record
	err := SetResponse(ctx, []byte{53})
	if err == nil {
		t.Errorf("Invalid record didn't return error.")
	}

	// Valid record
	_ = SetResponse(ctx, record)

	if !reflect.DeepEqual(expectedStatusCode, ctx.Response().StatusCode()) {
		t.Errorf("Expected status code: %v, Actual status code: %v", expectedStatusCode, ctx.Response().StatusCode())
	}

	if !reflect.DeepEqual(expectedResponseHeaders, ctx.GetRespHeaders()) {
		t.Errorf("Expected response headers: %v, Actual response headers: %v", expectedResponseHeaders, ctx.GetRespHeaders())
	}

	var responseBody map[string]interface{}
	json.Unmarshal(ctx.Response().Body(), &responseBody)

	if !reflect.DeepEqual(expectedResponseBody, responseBody) {
		t.Errorf("Expected response body: %v, Actual response body: %v", expectedResponseBody, responseBody)
	}
}
