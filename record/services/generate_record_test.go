package services

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/config"
	"wannabe/record/entities"
)

var testConfig = config.Records{
	Headers: config.HeadersToExclude{
		Exclude: []string{},
	},
}

var payload = entities.GenerateRecordPayload{
	Hash:       "testHash",
	Curl:       "testCurl",
	HttpMethod: "POST",
	Host:       "test.com",
	Path:       "/test",
	Query: map[string]string{
		"test": "test",
	},
	RequestHeaders: map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	},
	RequestBody: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
	StatusCode:  200,
	ResponseHeaders: map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	},
	ResponseBody: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
}

func TestGenerateRecord(t *testing.T) {
	expected := entities.Record{
		Request: entities.Request{
			Hash:       "testHash",
			Curl:       "testCurl",
			HttpMethod: "POST",
			Host:       "test.com",
			Path:       "/test",
			Query: map[string]string{
				"test": "test",
			},
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		},
		Response: entities.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		},
	}

	encodedRecord, _ := GenerateRecord(testConfig, payload)

	var record entities.Record

	_ = json.Unmarshal(encodedRecord, &record)

	if !reflect.DeepEqual(expected.Request, record.Request) {
		t.Errorf("Expected record request: %v, Actual record request: %v", expected, record)
	}

	if !reflect.DeepEqual(expected.Response, record.Response) {
		t.Errorf("Expected record response: %v, Actual record response: %v", expected, record)
	}
}
