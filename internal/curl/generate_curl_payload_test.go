package curl

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestGenerateCurlPayload(t *testing.T) {
	t.Run("generate curl payload", func(t *testing.T) {
		request := generateTestRequest()

		want := CurlPayload{
			Scheme:     "https",
			HttpMethod: "POST",
			Host:       "test.com",
			Path:       "/test",
			Query: map[string][]string{
				"test": {"test"},
			},
			RequestHeaders: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			// {"test":"test"}
			RequestBody: []byte{123, 34, 116, 101, 115, 116, 34, 58, 34, 116, 101, 115, 116, 34, 125},
		}

		got, _ := GenerateCurlPayload(request)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GenerateCurlPayload() = %v, want %v", got, want)
		}
	})
}

func generateTestRequest() *http.Request {
	httpMethod := "POST"
	url := "https://test.com/test?test=test"
	body := "{\"test\":\"test\"}"
	bodyBuffer := bytes.NewBufferString(body)

	request, _ := http.NewRequest(httpMethod, url, bodyBuffer)

	request.Header.Set("Accept", "test")
	request.Header.Set("Content-Type", "application/json")

	return request
}
