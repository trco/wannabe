package services

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
		fmt.Fprintln(w, "Test")
	}))
	defer testServer.Close()

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	// invalid testServer.URL
	err := FetchResponse(ctx, "http", "invalidUrl")
	if err == nil {
		t.Errorf("invalid url didn't return error.")
	}

	// valid testServer.URL
	_ = FetchResponse(ctx, "http", testServer.URL)

	expectedStatusCode := 200
	expectedResponseHeaderContentType := "application/json"
	expectedResponseHeaderAccept := "test"
	expectedResponseBody := []byte{84, 101, 115, 116, 10}

	if !reflect.DeepEqual(expectedStatusCode, ctx.Response().StatusCode()) {
		t.Errorf("expected status code: %v, actual status code: %v", expectedStatusCode, ctx.Response().StatusCode())
	}

	if !reflect.DeepEqual(expectedResponseHeaderContentType, ctx.GetRespHeader("Content-Type")) {
		t.Errorf("expected response header Content-Type: %v, actual response header Content-Type: %v", expectedResponseHeaderContentType, ctx.GetRespHeader("Content-Type"))
	}

	if !reflect.DeepEqual(expectedResponseHeaderAccept, ctx.GetRespHeader("Accept")) {
		t.Errorf("expected response header Accept: %v, actual response header Accept: %v", expectedResponseHeaderAccept, ctx.GetRespHeader("Accept"))
	}

	if !reflect.DeepEqual(expectedResponseBody, ctx.Response().Body()) {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, string(ctx.Response().Body()))
	}
}
