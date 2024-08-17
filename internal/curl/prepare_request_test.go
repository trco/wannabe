package curl

import (
	"io"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	tests := []struct {
		name       string
		httpMethod string
		url        string
		body       string
	}{
		{
			name:       "prepare POST request with body",
			httpMethod: "POST",
			url:        "http://test.com/test?test=test",
			body:       "{\"test\":\"test\"}",
		},
		{
			name:       "prepare GET request without body",
			httpMethod: "GET",
			url:        "http://test.com/test?test=test",
			body:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := PrepareRequest(tt.httpMethod, tt.url, tt.body)

			if request.Method != tt.httpMethod {
				t.Errorf("got http method %v, want http method %v", request.Method, tt.httpMethod)
			}

			if request.URL.String() != tt.url {
				t.Errorf("got url %v, want url %v", request.URL, tt.url)
			}

			if tt.httpMethod != "GET" {
				requestBody, _ := io.ReadAll(request.Body)
				if string(requestBody) != tt.body {
					t.Errorf("got body %v, want body %v", string(requestBody), tt.body)
				}
			}

		})
	}
}
