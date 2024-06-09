package actions

import (
	"io"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateRequest(t *testing.T) {
	recordRequest := types.Request{
		Hash:       "testHash",
		Curl:       "testCurl",
		HttpMethod: "POST",
		Host:       "test.com",
		Path:       "",
		Query: map[string][]string{
			"test": {"test"},
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"test"},
		},
		// {"test":"test"}
		Body: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
	}

	request, _ := GenerateRequest(recordRequest)

	expectedHttpMethod := "POST"

	if expectedHttpMethod != request.Method {
		t.Errorf("expected http method: %v, actual http method: %v", expectedHttpMethod, request.Method)
	}

	expectedHost := "test.com"

	if expectedHost != request.URL.Host {
		t.Errorf("expected host: %v, actual host: %v", expectedHost, request.URL.Host)
	}

	expectedPath := "/"

	if expectedPath != request.URL.Path {
		t.Errorf("expected path: %v, actual path: %v", expectedPath, request.URL.Path)
	}

	expectedQuery := "test=test"

	if expectedQuery != request.URL.RawQuery {
		t.Errorf("expected query: %v, actual query: %v", expectedQuery, request.URL.RawQuery)
	}

	expectedHeader := "application/json"

	if expectedHeader != request.Header.Get("Content-Type") {
		t.Errorf("expected Content-Type header: %v, actual Content-Type header: %v", expectedHeader, request.Header.Get("Content-Type"))
	}

	expectedRequestBody := []byte{34, 101, 119, 111, 103, 73, 67, 65, 103, 73, 110, 82, 108, 99, 51, 81, 105, 79, 105, 65, 105, 100, 71, 86, 122, 100, 67, 73, 75, 102, 81, 61, 61, 34}

	requestBody, _ := io.ReadAll(request.Body)

	if !reflect.DeepEqual(expectedRequestBody, requestBody) {
		t.Errorf("expected request body: %v, actual request body: %v", expectedRequestBody, requestBody)
	}
}
