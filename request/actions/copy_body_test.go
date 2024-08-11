package actions

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

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

			copiedBody, err := CopyBody(req)
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
