package actions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestFetchResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Accept", "test")
		w.Header().Set("Content-Length", "5")
		fmt.Fprintln(w, "Test")
	}))
	defer testServer.Close()

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	// Invalid testServer.URL
	err := FetchResponse(ctx, "invalidUrl")
	if err == nil {
		t.Errorf("Invalid url didn't return error.")
	}

	// Valid testServer.URL
	_ = FetchResponse(ctx, testServer.URL)

	expectedStatusCode := 200
	expectedResponseHeaderContentType := "application/json"
	expectedResponseHeaderAccept := "test"
	expectedResponseBody := []byte{84, 101, 115, 116, 10}

	if !reflect.DeepEqual(expectedStatusCode, ctx.Response().StatusCode()) {
		t.Errorf("Expected status code: %v, Actual status code: %v", expectedStatusCode, ctx.Response().StatusCode())
	}

	if !reflect.DeepEqual(expectedResponseHeaderContentType, ctx.GetRespHeader("Content-Type")) {
		t.Errorf("Expected response header Content-Type: %v, Actual response header Content-Type: %v", expectedResponseHeaderContentType, ctx.GetRespHeader("Content-Type"))
	}

	if !reflect.DeepEqual(expectedResponseHeaderAccept, ctx.GetRespHeader("Accept")) {
		t.Errorf("Expected response header Accept: %v, Actual response header Accept: %v", expectedResponseHeaderAccept, ctx.GetRespHeader("Accept"))
	}

	if !reflect.DeepEqual(expectedResponseBody, ctx.Response().Body()) {
		t.Errorf("Expected response body: %v, Actual response body: %v", expectedResponseBody, string(ctx.Response().Body()))
	}
}
