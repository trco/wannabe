package actions

import (
	"net/http"
	"testing"
	"wannabe/types"
)

func TestSetHeaders(t *testing.T) {
	request := generateTestRequest()

	headers := []types.Header{
		{Key: "Authorization", Value: "test access token"},
		{Key: "X-test-header", Value: "test value"},
	}

	SetHeaders(request, headers)

	if request.Header["Authorization"][0] != headers[0].Value {
		t.Errorf("expected Authorization header: %s, actual Authorization header: %s", request.Header["Authorization"][0], headers[0].Value)
	}

	if request.Header[http.CanonicalHeaderKey("X-test-header")][0] != headers[1].Value {
		t.Errorf("expected X-test-header header: %s, actual X-test-header header: %s", request.Header[http.CanonicalHeaderKey("X-test-header")][0], headers[1].Value)
	}
}
