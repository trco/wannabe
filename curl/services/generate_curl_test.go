package services

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"wannabe/types"
)

var wannabeA = types.Wannabe{
	RequestMatching: types.RequestMatching{
		Headers: types.Headers{
			Include: []string{"Content-Type", "Accept"},
		},
	},
}

func TestGenerateCurl(t *testing.T) {
	setHeaders(originalRequest)

	expcetedCurl := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'test.com/test?test=test'"

	curl, _ := GenerateCurl(originalRequest, wannabeA)

	if !reflect.DeepEqual(expcetedCurl, curl) {
		t.Errorf("expected curl: %v, actual curl: %v", expcetedCurl, curl)
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
