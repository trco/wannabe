package actions

import (
	"io"
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
		Body: map[string]interface{}{
			"test": "test",
		},
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

	expectedRequestBody := "{\"test\":\"test\"}"

	body, _ := io.ReadAll(request.Body)
	requestBody := string(body)

	if expectedRequestBody != requestBody {
		t.Errorf("expected request body: %v, actual request body: %v", expectedRequestBody, requestBody)
	}
}
