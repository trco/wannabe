package record

import (
	"io"
	"testing"
)

func TestGenerateRequest(t *testing.T) {
	t.Run("generate request", func(t *testing.T) {
		recordRequest := Request{
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

		got, _ := GenerateRequest(recordRequest)

		wantHttpMethod := "POST"
		wantHost := "test.com"
		wantPath := "/"
		wantQuery := "test=test"
		wantHeader := "application/json"
		wantRequestBody := "{\"test\":\"test\"}"

		body, _ := io.ReadAll(got.Body)
		gotRequestBody := string(body)

		if got.Method != wantHttpMethod {
			t.Errorf("got http method = %v, want http method %v", got.Method, wantHttpMethod)
		}

		if got.URL.Host != wantHost {
			t.Errorf("got host = %v, want host %v", got.URL.Host, wantHost)
		}

		if got.URL.Path != wantPath {
			t.Errorf("got path = %v, want path %v", got.URL.Path, wantPath)
		}

		if got.URL.RawQuery != wantQuery {
			t.Errorf("got query = %v, want query %v", got.URL.RawQuery, wantQuery)
		}

		if got.Header.Get("Content-Type") != wantHeader {
			t.Errorf("got Content-Type header = %v, want Content-Type header %v", got.Header.Get("Content-Type"), wantHeader)
		}

		if gotRequestBody != wantRequestBody {
			t.Errorf("got request body = %v, want request body %v", gotRequestBody, wantRequestBody)
		}
	})
}
