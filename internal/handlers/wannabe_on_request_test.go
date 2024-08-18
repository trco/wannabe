package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
	"github.com/trco/wannabe/internal/record"
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

func TestSetResponse(t *testing.T) {
	tests := []struct {
		name                string
		record              []byte
		request             *http.Request
		wantStatusCode      int
		wantResponseHeaders map[string][]string
		wantResponseBody    map[string]interface{}
		wantErr             string
	}{
		{
			name:                "invalid record",
			record:              []byte{53},
			request:             request,
			wantStatusCode:      0,
			wantResponseHeaders: nil,
			wantResponseBody:    nil,
			wantErr:             "SetResponse: failed unmarshaling record: json: cannot unmarshal number into Go value of type Record",
		},
		{
			name:           "valid record",
			record:         testRecord,
			request:        request,
			wantStatusCode: 200,
			wantResponseHeaders: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			wantResponseBody: map[string]interface{}{
				"test": "test",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setResponse(tt.record, tt.request)

			if err != nil {
				if err.Error() != tt.wantErr {
					t.Errorf("SetResponse() error = %v, wantErr %v", err.Error(), tt.wantErr)
				}
				return
			}

			if !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("got status code %v, want status code %v", got.StatusCode, tt.wantStatusCode)
			}

			if !(len(got.Header) == len(tt.wantResponseHeaders)) {
				t.Errorf("got %v response headers , want %v response headers", len(got.Header), len(tt.wantResponseHeaders))
			}

			responseBodyEncoded, _ := io.ReadAll(got.Body)
			defer got.Body.Close()

			var gotResponseBody map[string]interface{}
			json.Unmarshal(responseBodyEncoded, &gotResponseBody)

			if !reflect.DeepEqual(gotResponseBody, tt.wantResponseBody) {
				t.Errorf("got response body: %v, want response body: %v", gotResponseBody, tt.wantResponseBody)
			}
		})
	}
}

var request = &http.Request{
	Method: "POST",
	Host:   "test.com",
	URL: &url.URL{
		Path:     "/test",
		RawQuery: "test=test",
	},
	Header: map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	},
	Body: io.NopCloser(bytes.NewBuffer([]byte(`{"test": "test"}`))),
}

var testRecord, _ = json.Marshal(record.Record{
	Request: record.Request{
		HttpMethod: "POST",
		Host:       "test.com",
		Path:       "/test",
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
	},
	Response: record.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"test"},
		},
		Body: map[string]interface{}{
			"test": "test",
		},
	},
})

func TestEncodeBody(t *testing.T) {
	tests := []struct {
		name            string
		decodedBody     interface{}
		contentType     []string
		contentEncoding []string
		want            []byte
		wantErr         bool
	}{
		{
			name:            "json content type",
			decodedBody:     map[string]string{"key": "value"},
			contentType:     []string{"application/json"},
			contentEncoding: []string{},
			want:            []byte(`{"key":"value"}`),
			wantErr:         false,
		},
		{
			name:            "xml content type",
			decodedBody:     map[string]interface{}{"root": map[string]interface{}{"key": "value"}},
			contentType:     []string{"application/xml"},
			contentEncoding: []string{},
			want:            nil,
			wantErr:         false,
		},
		{
			name:            "text/plain content type",
			decodedBody:     "plain text",
			contentType:     []string{"text/plain"},
			contentEncoding: []string{},
			want:            []byte("plain text"),
			wantErr:         false,
		},
		{
			name:            "text/html content type",
			decodedBody:     "\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e",
			contentType:     []string{"text/html"},
			contentEncoding: []string{},
			want:            []byte("\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e"),
			wantErr:         false,
		},
		{
			name:            "unsupported content type",
			decodedBody:     "unsupported",
			contentType:     []string{"unsupported"},
			contentEncoding: []string{},
			want:            nil,
			wantErr:         true,
		},
		{
			name:            "json content type, gzip content encoding",
			decodedBody:     map[string]string{"key": "value"},
			contentType:     []string{"application/json"},
			contentEncoding: []string{"gzip"},
			want:            getCompressedEncodedBody([]byte(`{"key":"value"}`)),
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeBody(tt.decodedBody, tt.contentType, tt.contentEncoding)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "xml content type" {
				mapValue := mxj.Map(tt.decodedBody.(map[string]interface{}))
				body, _ := mapValue.Xml()

				tt.want = body
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCompressedEncodedBody(data []byte) []byte {
	compressedEncodedBody, _ := record.Gzip(data)
	return compressedEncodedBody
}

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
