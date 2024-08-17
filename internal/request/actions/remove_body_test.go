package actions

import (
	"net/http"
	"testing"
)

func TestRemoveBody(t *testing.T) {
	t.Run("remove body", func(t *testing.T) {
		httpMethod := "GET"
		url := "http://test.com"

		request, _ := http.NewRequest(httpMethod, url, nil)
		request.Header.Set("Content-Length", "4")
		request.Header.Set("Content-Type", "application/json")

		got := RemoveBody(request)

		if got.Body != http.NoBody {
			t.Errorf("RemoveBody() = body %v, body %v", got.Body, http.NoBody)
		}

		if len(got.Header) != 0 {
			t.Errorf("RemoveBody() = headers count %v, want %v", len(got.Header), 0)
		}
	})
}
