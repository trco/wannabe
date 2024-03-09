package services

import (
	"reflect"
	"testing"
	"wannabe/config"
)

var testConfigA = config.Config{
	Server: "https:/test.com",
	RequestMatching: config.RequestMatching{
		Headers: config.Headers{
			Include: []string{"Content-Type", "Accept"},
		},
	},
}

func TestGenerateCurl(t *testing.T) {
	method := ("POST")
	path := "test"
	queries := map[string]string{
		"test": "test",
	}
	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	}
	// {"test":"test"}
	body := []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125}

	expcetedCurl := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'https:/test.com/test?test=test'"

	curl, _ := GenerateCurl(method, path, queries, headers, body, testConfigA)

	if !reflect.DeepEqual(expcetedCurl, curl) {
		t.Errorf("Expected curl: %v, Actual curl: %v", expcetedCurl, curl)
	}
}
