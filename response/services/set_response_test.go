package services

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/config"
	"wannabe/providers"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var testConfigB = config.Config{
	StorageProvider: config.StorageProvider{
		Type: "filesystem",
		FilesystemConfig: config.FilesystemConfig{
			Folder: "/var/folders/6z/9bvblj5j2s9bngjcnr18jls80000gn/T",
			Format: "json",
		},
	},
}

var filesystemProviderB = providers.FilesystemProvider{
	Config: testConfigB.StorageProvider,
}

func TestSetResponse(t *testing.T) {
	testHash := "testHash3"
	// testHashInvalid := "testHashInvalid"

	testRecord, _ := json.Marshal(entities.Record{
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

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	expectedStatusCode := 200
	expectedResponseHeaders := map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	}
	expectedResponseBody := map[string]interface{}{
		"test": "test",
	}

	// FIXME update tests
	// add record
	filesystemProviderB.InsertRecords([]string{testHash}, [][]byte{testRecord})

	// invalid hash without record
	err := SetResponse(ctx, testRecord)

	if err == nil {
		t.Errorf("Invalid hash didn't return error.")
	}

	// valid hash with record
	_ = SetResponse(ctx, testRecord)

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
