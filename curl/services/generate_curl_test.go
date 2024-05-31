package services

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"wannabe/types"
)

var wannabe = types.Wannabe{
	RequestMatching: types.RequestMatching{
		Headers: types.Headers{
			Include: []string{"Content-Type", "Accept"},
		},
	},
}

func TestGenerateCurl(t *testing.T) {
	request := generateTestRequest()

	expectedCurl := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'test.com/test?test=test'"

	curl, _ := GenerateCurl(request, wannabe)

	if !reflect.DeepEqual(expectedCurl, curl) {
		t.Errorf("expected curl: %v, actual curl: %v", expectedCurl, curl)
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
