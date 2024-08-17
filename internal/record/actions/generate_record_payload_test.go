package actions

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/trco/wannabe/types"
)

func TestGenerateRecordPayload(t *testing.T) {
	t.Run("generate record payload", func(t *testing.T) {
		request := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "https",
				Host:   "test.com",
				Path:   "/test",
			},
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBufferString("request body")),
		}

		response := &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(bytes.NewBufferString("response body")),
		}

		hash := "test hash"
		curl := "test curl"

		session := &MockSession{
			props: map[string]interface{}{
				"hash": hash,
				"curl": curl,
			},
			req: request,
			res: response,
		}

		want := types.RecordPayload{
			Hash:            hash,
			Curl:            curl,
			Scheme:          request.URL.Scheme,
			HttpMethod:      request.Method,
			Host:            request.URL.Host,
			Path:            request.URL.Path,
			Query:           request.URL.Query(),
			RequestHeaders:  request.Header,
			RequestBody:     []byte{114, 101, 113, 117, 101, 115, 116, 32, 98, 111, 100, 121},
			StatusCode:      response.StatusCode,
			ResponseHeaders: response.Header,
			ResponseBody:    []byte{114, 101, 115, 112, 111, 110, 115, 101, 32, 98, 111, 100, 121},
		}

		got, _ := GenerateRecordPayload(session, hash, curl)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GenerateRecordPayload() = %v, want %v", got, wantRecordTwo)
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
