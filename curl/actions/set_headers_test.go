package actions

import (
	"bytes"
	"net/http"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	// prepare request
	httpMethod := "POST"
	url := "testUrl"
	bodyBuffer := bytes.NewBufferString("")
	request, _ := http.NewRequest(httpMethod, url, bodyBuffer)
	headers := []Header{
		{Key: "Accept", Value: "test1,test2,test3"},
		{Key: "Authorization", Value: "test access token"},
		{Key: "Content-Type", Value: "application/json"},
		{Key: "X-test-header", Value: "test value"},
	}

	// test setting of headers
	SetHeaders(request, headers)

	if request.Header["Accept"][0] != headers[0].Value {
		t.Errorf("Expected Accept header: %s, Actual Accept header: %s", request.Header["Accept"][0], headers[0].Value)
	}

	if request.Header["Authorization"][0] != headers[1].Value {
		t.Errorf("Expected Authorization header: %s, Actual Authorization header: %s", request.Header["Authorization"][0], headers[1].Value)
	}

	if request.Header["Content-Type"][0] != headers[2].Value {
		t.Errorf("Expected Content-Type header: %s, Actual Content-Type header: %s", request.Header["Content-Type"][0], headers[2].Value)
	}

	if request.Header[http.CanonicalHeaderKey("X-test-header")][0] != headers[3].Value {
		t.Errorf("Expected X-test-header header: %s, Actual X-test-header header: %s", request.Header[http.CanonicalHeaderKey("X-test-header")][0], headers[3].Value)
	}
}
