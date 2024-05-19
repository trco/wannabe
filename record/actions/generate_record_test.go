package actions

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/record/entities"
	"wannabe/types"
)

func TestGenerateRecord(t *testing.T) {
	encodedRecord, _ := GenerateRecord(testConfigA, payload)

	var record entities.Record

	_ = json.Unmarshal(encodedRecord, &record)

	if !reflect.DeepEqual(expectedRecordC.Request, record.Request) {
		t.Errorf("expected record request: %v, actual record request: %v", expectedRecordC, record)
	}

	if !reflect.DeepEqual(expectedRecordC.Response, record.Response) {
		t.Errorf("expected record response: %v, actual record response: %v", expectedRecordC, record)
	}
}

// reusable variables
var testConfigA = types.Records{
	Headers: types.HeadersToExclude{
		Exclude: []string{"Header-To-Exclude"},
	},
}

var payload = entities.RecordPayload{
	Hash:       "testHash",
	Curl:       "testCurl",
	HttpMethod: "POST",
	Host:       "test.com",
	Path:       "/test",
	Query: map[string][]string{
		"test": {"test"},
	},
	RequestHeaders: map[string][]string{
		"Content-Type":      {"application/json"},
		"Accept":            {"test"},
		"Header-To-Exclude": {"test"},
	},
	// {"test":"test"}
	RequestBody: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
	StatusCode:  200,
	ResponseHeaders: map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	},
	// {"test":"test"}
	ResponseBody: []byte{123, 10, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 125},
}

var expectedRecordC = entities.Record{
	Request: entities.Request{
		Hash:       "testHash",
		Curl:       "testCurl",
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
