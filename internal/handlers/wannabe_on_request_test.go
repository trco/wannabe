package handlers

import (
	"bytes"
	"io"
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

		got := processRequest(request)
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

func TestRemoveBody(t *testing.T) {
	t.Run("remove body", func(t *testing.T) {
		httpMethod := "GET"
		url := "http://test.com"

		request, _ := http.NewRequest(httpMethod, url, nil)
		request.Header.Set("Content-Length", "4")
		request.Header.Set("Content-Type", "application/json")

		got := removeBody(request)

		if got.Body != http.NoBody {
			t.Errorf("RemoveBody() = body %v, body %v", got.Body, http.NoBody)
		}

		if len(got.Header) != 0 {
			t.Errorf("RemoveBody() = headers count %v, want %v", len(got.Header), 0)
		}
	})
}

func TestProcessScheme(t *testing.T) {
	t.Run("process scheme", func(t *testing.T) {
		httpMethod := "GET"
		url := "http://test.com"

		request, _ := http.NewRequest(httpMethod, url, nil)
		request.RequestURI = "https://test.com"

		processedRequest := processScheme(request)

		got := processedRequest.URL.Scheme
		want := "https"

		if got != want {
			t.Errorf("ProcessScheme() = scheme %v, want %v", got, want)
		}
	})
}

func TestCopyBody(t *testing.T) {
	tests := []struct {
		name          string
		httpMethod    string
		url           string
		body          string
		expectedError bool
	}{
		{
			name:          "Non-empty body",
			httpMethod:    "POST",
			url:           "http://test.com",
			body:          "{\"test\":\"test\"}",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.httpMethod, tt.url, bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			copiedBody, err := copyBody(req)
			if (err != nil) != tt.expectedError {
				t.Fatalf("CopyBody() error = %v, expectedError %v", err, tt.expectedError)
			}

			copiedBodyBytes, err := io.ReadAll(copiedBody)
			if err != nil {
				t.Fatalf("Failed to read copied body: %v", err)
			}

			if string(copiedBodyBytes) != tt.body {
				t.Errorf("CopyBody() = %s, want %s", string(copiedBodyBytes), tt.body)
			}
		})
	}
}
