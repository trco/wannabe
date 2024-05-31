package actions

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateCurlPayload(t *testing.T) {
	request := generateTestRequest()

	expcetedPayload := types.CurlPayload{
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

	payload, _ := GenerateCurlPayload(request)

	if !reflect.DeepEqual(expcetedPayload, payload) {
		t.Errorf("expected payload: %v, actual payload: %v", expcetedPayload, payload)
	}
}

func generateTestRequest() *http.Request {
	httpMethod := "POST"
	url := "http://test.com/test?test=test"
	body := "{\"test\":\"test\"}"
	bodyBuffer := bytes.NewBufferString(body)

	request, _ := http.NewRequest(httpMethod, url, bodyBuffer)

	request.Header.Set("Accept", "test")
	request.Header.Set("Content-Type", "application/json")

	return request
}
