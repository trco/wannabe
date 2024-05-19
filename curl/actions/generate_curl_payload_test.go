package actions

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateCurlPayload(t *testing.T) {
	setHeaders(originalRequest)

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

	payload, _ := GenerateCurlPayload(originalRequest)

	if !reflect.DeepEqual(expcetedPayload, payload) {
		t.Errorf("expected payload: %v, actual payload: %v", expcetedPayload, payload)
	}
}

// reusable variables and methods
var requestBody = "{\"test\":\"test\"}"
var bodyBuffer = bytes.NewBufferString(requestBody)
var originalRequest, _ = http.NewRequest("POST", "http://test.com/test?test=test", bodyBuffer)

func setHeaders(originalRequest *http.Request) {
	originalRequest.Header.Set("Accept", "test")
	originalRequest.Header.Set("Content-Type", "application/json")
}
