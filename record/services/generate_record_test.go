package services

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/config"
	"wannabe/record/common"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var testConfig = config.Records{
	Headers: config.HeadersToExclude{
		Exclude: []string{},
	},
}

var testServer = "test.com"
var testCurl = "testCurl"
var testHash = "testHash"

func TestGenerateRecord(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

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

	expected, _ := json.Marshal(common.Record{
		Request: common.Request{
			Hash:       "testHash",
			Curl:       "testCurl",
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
		Response: common.Response{
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

	record, _ := GenerateRecord(ctx, testConfig, testServer, testCurl, testHash)

	if !reflect.DeepEqual(expected, record) {
		t.Errorf("Expected record: %v, Actual record: %v", expected, record)
	}
}
