package actions

import (
	"testing"
)

func TestSetResponse(t *testing.T) {
	// app := fiber.New()
	// ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	// record, _ := json.Marshal(entities.Record{
	// 	Request: entities.Request{
	// 		HttpMethod: "POST",
	// 		Host:       "test.com",
	// 		Path:       "/test",
	// 		Query: map[string]string{
	// 			"test": "test",
	// 		},
	// 		Headers: map[string][]string{
	// 			"Content-Type": {"application/json"},
	// 			"Accept":       {"test"},
	// 		},
	// 		Body: map[string]interface{}{
	// 			"test": "test",
	// 		},
	// 	},
	// 	Response: entities.Response{
	// 		StatusCode: 200,
	// 		Headers: map[string][]string{
	// 			"Content-Type": {"application/json"},
	// 			"Accept":       {"test"},
	// 		},
	// 		Body: map[string]interface{}{
	// 			"test": "test",
	// 		},
	// 	},
	// })

	// expectedStatusCode := 200
	// expectedResponseHeaders := map[string][]string{
	// 	"Content-Type": {"application/json"},
	// 	"Accept":       {"test"},
	// }
	// expectedResponseBody := map[string]interface{}{
	// 	"test": "test",
	// }

	// // invalid record
	// err := SetResponse(ctx, []byte{53})
	// if err == nil {
	// 	t.Errorf("invalid record didn't return error.")
	// }

	// // valid record
	// _ = SetResponse(ctx, record)

	// if !reflect.DeepEqual(expectedStatusCode, ctx.Response().StatusCode()) {
	// 	t.Errorf("expected status code: %v, actual status code: %v", expectedStatusCode, ctx.Response().StatusCode())
	// }

	// if !reflect.DeepEqual(expectedResponseHeaders, ctx.GetRespHeaders()) {
	// 	t.Errorf("expected response headers: %v, actual response headers: %v", expectedResponseHeaders, ctx.GetRespHeaders())
	// }

	// var responseBody map[string]interface{}
	// json.Unmarshal(ctx.Response().Body(), &responseBody)

	// if !reflect.DeepEqual(expectedResponseBody, responseBody) {
	// 	t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	// }
}
