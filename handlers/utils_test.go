package handlers

import (
	"testing"
)

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}

// FIXME
// func TestInternalError(t *testing.T) {
// 	app := fiber.New()
// 	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

// 	testError := &TestError{"test error"}

// 	_ = internalError(ctx, testError)

// 	expectedStatusCode := 500
// 	expectedResponseBody := "{\"statusCode\":500,\"error\":\"test error\"}"

// 	if !reflect.DeepEqual(expectedStatusCode, ctx.Response().StatusCode()) {
// 		t.Errorf("expected status code: %v, actual status code: %v", expectedStatusCode, ctx.Response().StatusCode())
// 	}

// 	if !reflect.DeepEqual(expectedResponseBody, string(ctx.Response().Body())) {
// 		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, string(ctx.Response().Body()))
// 	}
// }

func TestCheckDuplicates(t *testing.T) {
	// duplicates exist
	slice := []string{"a", "b", "c"}

	duplicatesExist := checkDuplicates(slice, "a")

	if duplicatesExist != true {
		t.Errorf("no duplicates detected although they are present in slice")
	}

	// no duplicates
	duplicatesExist = checkDuplicates(slice, "d")

	if duplicatesExist == true {
		t.Errorf("duplicates detected although they are not present")
	}
}

func TestProcessRecordValidation(t *testing.T) {
	count := 0
	var recordProcessingDetails []RecordProcessingDetails

	processRecordValidation(&recordProcessingDetails, "test hash", "test message", &count)

	if recordProcessingDetails[0].Hash != "test hash" || recordProcessingDetails[0].Message != "test message" || count != 1 {
		t.Errorf("record processing details not valid, expected: hash: 'test hash', message: 'test message', count: 1, actual: hash: '%v', message: '%v', count: %v", recordProcessingDetails[0].Hash, recordProcessingDetails[0].Message, count)
	}
}
