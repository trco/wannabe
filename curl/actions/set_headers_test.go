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

	t.Run("set headers", func(t *testing.T) {
		SetHeaders(request, headers)

		if request.Header["Authorization"][0] != headers[0].Value {
			t.Errorf("SetHeaders() got Authorization header %v, want Authorization header %v", request.Header["Authorization"][0], headers[0].Value)
		}

		if request.Header[http.CanonicalHeaderKey("X-test-header")][0] != headers[1].Value {
			t.Errorf("SetHeaders() got X-test-header header %v, want X-test-header header  %s", request.Header[http.CanonicalHeaderKey("X-test-header")][0], headers[1].Value)
		}
	})

}
