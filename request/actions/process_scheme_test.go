package actions

import (
	"net/http"
	"testing"
)

func TestProcessScheme(t *testing.T) {
	t.Run("process https", func(t *testing.T) {
		httpMethod := "GET"
		url := "http://test.com"

		request, _ := http.NewRequest(httpMethod, url, nil)
		request.RequestURI = "https://test.com"

		processedRequest := ProcessScheme(request)

		got := processedRequest.URL.Scheme
		want := "https"

		if got != want {
			t.Errorf("ProcessScheme() = scheme %v, want %v", got, want)
		}
	})
}
