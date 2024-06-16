package actions

import (
	"io"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	t.Run("prepare request", func(t *testing.T) {
		httpMethod := "POST"
		url := "http://test.com/test?test=test"
		body := "{\"test\":\"test\"}"

		request, _ := PrepareRequest(httpMethod, url, body)
		requestBody, _ := io.ReadAll(request.Body)

		if request.Method != httpMethod {
			t.Errorf("got http method %v, want http method %v", httpMethod, request.Method)
		}

		if request.URL.String() != url {
			t.Errorf("got url %v, want url %v", url, request.URL)
		}

		if string(requestBody) != body {
			t.Errorf("got body %v, want body %v", body, string(requestBody))
		}
	})
}
