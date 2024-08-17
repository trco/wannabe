package services

import (
	"net/http"
	"testing"
)

func TestProcessRequest(t *testing.T) {
	t.Run("process request", func(t *testing.T) {
		httpMethod := "GET"
		url := "http://test.com"

		request, _ := http.NewRequest(httpMethod, url, nil)
		request.Header.Set("Content-Length", "4")
		request.Header.Set("Content-Type", "application/json")
		request.RequestURI = "https://test.com"

		got := ProcessRequest(request)
		want := "https"

		if got.Body != http.NoBody {
			t.Errorf("ProcessRequest() = body %v, body %v", got.Body, http.NoBody)
		}

		if len(got.Header) != 0 {
			t.Errorf("ProcessRequest() = headers count %v, want %v", len(got.Header), 0)
		}

		if got.URL.Scheme != want {
			t.Errorf("ProcessRequest() = scheme %v, want %v", got, want)
		}
	})
}
