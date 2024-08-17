package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trco/wannabe/types"
)

func TestInternalErrorOnRequest(t *testing.T) {
	t.Run("internal error on request", func(t *testing.T) {
		session := &MockSession{
			props: make(map[string]interface{}),
		}
		request, _ := http.NewRequest("GET", "test.com", nil)

		_, response := InternalErrorOnRequest(session, request, &TestError{"test error"})
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

func TestInternalErrorOnResponse(t *testing.T) {
	t.Run("internal error on request", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://test.com", nil)

		response := InternalErrorOnResponse(request, &TestError{"test error"})
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

func TestPrepareResponseBody(t *testing.T) {
	t.Run("prepare response body", func(t *testing.T) {
		reader := PrepareResponseBody(&TestError{"test error"})
		want := "{\"error\":\"test error\"}"

		got, _ := io.ReadAll(reader)

		if !bytes.Equal([]byte(want), got) {
			t.Errorf("PrepareResponseBody() = %v, want %v", string(got), want)
		}
	})
}

func TestShouldSkipResponseProcessing(t *testing.T) {
	tests := []struct {
		name    string
		session *MockSession
		want    bool
	}{
		{
			name:    "'blocked' and 'responseSetFromRecord' not set on session",
			session: &MockSession{props: make(map[string]interface{})},
			want:    false,
		},
		{
			name:    "'blocked' set to 'true' on session",
			session: &MockSession{props: map[string]interface{}{"blocked": true}},
			want:    true,
		},
		{
			name:    "'responseSetFromRecord' set to 'true' on session",
			session: &MockSession{props: map[string]interface{}{"responseSetFromRecord": true}},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldSkipResponseProcessing(tt.session)

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
			hash, curl, _ := GetHashAndCurlFromSession(tt.session)

			if hash != tt.wantHash || curl != tt.wantCurl {
				t.Errorf("GetHashAndCurlFromSession() = hash %v, curl %v, want hash %v, want curl %v", hash, curl, tt.wantHash, tt.wantCurl)
			}
		})
	}
}

func TestInternalErrorApi(t *testing.T) {
	t.Run("internal error api", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := &TestError{"test error"}
		status := 500

		InternalErrorApi(w, err, status)
		responseBody := w.Body.String()

		wantStatusCode := http.StatusInternalServerError
		wantContentType := "application/json"
		wantResponseBody := "{\"error\":\"test error\"}\n"

		if w.Code != wantStatusCode {
			t.Errorf("got status code %v, want status code %v", w.Code, wantStatusCode)
		}

		if w.Header().Get("Content-Type") != wantContentType {
			t.Errorf("got content type %v, want content type %v", w.Header().Get("Content-Type"), wantContentType)
		}

		if responseBody != wantResponseBody {
			t.Errorf("got response body %v, want response body %v", responseBody, wantResponseBody)
		}
	})
}

func TestApiResponse(t *testing.T) {
	t.Run("api response", func(t *testing.T) {
		w := httptest.NewRecorder()
		response := map[string]string{"key": "value"}

		ApiResponse(w, response)
		responseBody := w.Body.String()

		wantStatusCode := http.StatusOK
		wantContentType := "application/json"
		wantResponseBody := "{\"key\":\"value\"}\n"

		if w.Code != wantStatusCode {
			t.Errorf("got status code %v, want status code %v", w.Code, wantStatusCode)
		}

		if w.Header().Get("Content-Type") != wantContentType {
			t.Errorf("got content type %v, want content type %v", w.Header().Get("Content-Type"), wantContentType)
		}

		if responseBody != wantResponseBody {
			t.Errorf("got response body %v, want response body %v", responseBody, wantResponseBody)
		}
	})
}

func TestCheckDuplicates(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "duplicates exist",
			slice: []string{"a", "b", "c"},
			item:  "a",
			want:  true,
		},
		{
			name:  "no duplicates",
			slice: []string{"a", "b", "c"},
			item:  "d",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckDuplicates(tt.slice, tt.item)

			if got != tt.want {
				t.Errorf("CheckDuplicates() = %v, got: %v", got, tt.want)
			}
		})
	}
}

func TestProcessRecordValidation(t *testing.T) {
	t.Run("process record validation", func(t *testing.T) {
		count := 0
		var recordProcessingDetails []types.RecordProcessingDetails

		ProcessRecordValidation(&recordProcessingDetails, "test hash", "test message", &count)

		wantHash := "test hash"
		wantMessage := "test message"
		wantCount := 1

		gotHash := recordProcessingDetails[0].Hash
		gotMessage := recordProcessingDetails[0].Message
		gotCount := 1

		if gotHash != wantHash {
			t.Errorf("got hash %v, want hash %v", gotHash, wantHash)
		}

		if gotMessage != wantMessage {
			t.Errorf("got message %v, want message %v", gotMessage, wantMessage)
		}

		if gotCount != wantCount {
			t.Errorf("got count %v, want count %v", gotCount, wantCount)
		}
	})
}

type MockSession struct {
	req   *http.Request
	res   *http.Response
	props map[string]interface{}
}

func (m *MockSession) SetProp(key string, value interface{}) {
	m.props[key] = value
}

func (m *MockSession) GetProp(key string) (interface{}, bool) {
	v, ok := m.props[key]
	return v, ok
}

// Request returns the HTTP request of this session
func (m *MockSession) Request() *http.Request {
	return m.req
}

// Response returns the HTTP response of this session
func (m *MockSession) Response() *http.Response {
	return m.res
}

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}
