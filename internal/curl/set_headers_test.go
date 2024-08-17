package curl

import (
	"net/http"
	"testing"
)

func TestSetHeaders(t *testing.T) {
	t.Run("set headers", func(t *testing.T) {
		request := generateTestRequest()

		headers := []header{
			{key: "Authorization", value: "test access token"},
			{key: "X-test-header", value: "test value"},
		}

		setHeaders(request, headers)

		if request.Header["Authorization"][0] != headers[0].value {
			t.Errorf("setHeaders() got Authorization header %v, want Authorization header %v", request.Header["Authorization"][0], headers[0].value)
		}

		if request.Header[http.CanonicalHeaderKey("X-test-header")][0] != headers[1].value {
			t.Errorf("setHeaders() got X-test-header header %v, want X-test-header header  %s", request.Header[http.CanonicalHeaderKey("X-test-header")][0], headers[1].value)
		}
	})

}
