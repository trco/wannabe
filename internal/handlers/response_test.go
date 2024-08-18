package handlers

import (
	"io"
	"net/http"
	"testing"
)

func TestShouldSkipResponseProcessing(t *testing.T) {
	tests := []struct {
		name    string
		session *MockSession
		want    bool
	}{
		{
			name:    "'blocked' and 'response' not set on session",
			session: &MockSession{props: make(map[string]interface{})},
			want:    false,
		},
		{
			name:    "'blocked' set to 'true' on session",
			session: &MockSession{props: map[string]interface{}{"blocked": true}},
			want:    true,
		},
		{
			name:    "'response' set to 'true' on session",
			session: &MockSession{props: map[string]interface{}{"response": true}},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipResponseProcessing(tt.session)

			if got != tt.want {
				t.Errorf("ShouldSkipResponseProcessing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHashAndCurlFromSession(t *testing.T) {
	tests := []struct {
		name     string
		session  *MockSession
		wantHash string
		wantCurl string
	}{
		{
			name: "no hash",
			session: &MockSession{
				props: map[string]interface{}{
					"noHash": "test hash",
				},
			},
			wantHash: "",
		},
		{
			name: "hash not string",
			session: &MockSession{
				props: map[string]interface{}{
					"hash": true,
				},
			},
			wantHash: "",
		},
		{
			name: "no curl",
			session: &MockSession{
				props: map[string]interface{}{
					"hash":   "test hash",
					"noCurl": "test curl",
				},
			},
			wantHash: "",
			wantCurl: "",
		},
		{
			name: "curl not string",
			session: &MockSession{
				props: map[string]interface{}{
					"hash": "test hash",
					"curl": true,
				},
			},
			wantHash: "",
			wantCurl: "",
		},
		{
			name: "no hash and curl",
			session: &MockSession{
				props: map[string]interface{}{
					"noHash": "test hash",
					"noCurl": "test curl",
				},
			},
			wantHash: "",
			wantCurl: "",
		},
		{
			name: "hash and curl",
			session: &MockSession{
				props: map[string]interface{}{
					"hash": "test hash",
					"curl": "test curl",
				},
			},
			wantHash: "test hash",
			wantCurl: "test curl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, curl, _ := getHashAndCurlFromSession(tt.session)

			if hash != tt.wantHash || curl != tt.wantCurl {
				t.Errorf("GetHashAndCurlFromSession() = hash %v, curl %v, want hash %v, want curl %v", hash, curl, tt.wantHash, tt.wantCurl)
			}
		})
	}
}

func TestInternalErrorOnResponse(t *testing.T) {
	t.Run("internal error on request", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://test.com", nil)

		response := internalErrorOnResponse(request, &TestError{"test error"})
		body, _ := io.ReadAll(response.Body)
		responseBody := string(body)

		wantContentType := "application/json"
		wantResponseBody := "{\"error\":\"test error\"}"
		wantStatusCode := http.StatusInternalServerError

		if response.StatusCode != wantStatusCode {
			t.Errorf("got status code %v, want status code %v", response.StatusCode, wantStatusCode)
		}

		if response.Header.Get("Content-Type") != wantContentType {
			t.Errorf("got content type %v, want content type %v", response.Header.Get("Content-Type"), wantContentType)
		}

		if responseBody != wantResponseBody {
			t.Errorf("got response body %v, want response body %v", responseBody, wantResponseBody)
		}
	})
}
