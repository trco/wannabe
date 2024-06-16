package services

import (
	"encoding/json"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateRecord(t *testing.T) {
	t.Run("generate record", func(t *testing.T) {
		encodedRecord, _ := GenerateRecord(testConfig, payload)

		var record types.Record

		_ = json.Unmarshal(encodedRecord, &record)

		if !reflect.DeepEqual(record.Request, wantRecordThree.Request) {
			t.Errorf("got record request: %v, want record request: %v", record.Request, wantRecordThree.Request)
		}

		if !reflect.DeepEqual(record.Response, wantRecordThree.Response) {
			t.Errorf("got record response: %v, want record response: %v", record.Response, wantRecordThree.Response)
		}

	})

}

var testConfig = types.Records{
	Headers: types.HeadersToRecord{
		Exclude: []string{"Header-To-Exclude"},
	},
}

var payload = types.RecordPayload{
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

var wantRecordThree = types.Record{
	Request: types.Request{
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
	Response: types.Response{
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
