package services

import (
	"reflect"
	"testing"
	"wannabe/config"
	"wannabe/curl/entities"
)

var testConfigA = config.Config{
	Server: "https:/test.com",
	RequestMatching: config.RequestMatching{
		Headers: config.Headers{
			Include: []string{"Content-Type", "Accept"},
		},
	},
}

var curlPayload = entities.GenerateCurlPayload{
	HttpMethod: "POST",
	Path:       "test",
	Query: map[string]string{
		"test": "test",
	},
	RequestHeaders: map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	},
	// {"test":"test"}
	RequestBody: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
}

func TestGenerateCurl(t *testing.T) {
	expcetedCurl := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'https:/test.com/test?test=test'"

	curl, _ := GenerateCurl(testConfigA, curlPayload)

	if !reflect.DeepEqual(expcetedCurl, curl) {
		t.Errorf("Expected curl: %v, Actual curl: %v", expcetedCurl, curl)
	}
}
