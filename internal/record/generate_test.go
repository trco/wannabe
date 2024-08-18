package record

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestGenerate(t *testing.T) {
	t.Run("generate record", func(t *testing.T) {
		encodedRecord, _ := Generate(testConfig, payload)

		var record Record

		_ = json.Unmarshal(encodedRecord, &record)

		if !reflect.DeepEqual(record.Request, wantRecordThree.Request) {
			t.Errorf("got record request: %v, want record request: %v", record.Request, wantRecordThree.Request)
		}

		if !reflect.DeepEqual(record.Response, wantRecordThree.Response) {
			t.Errorf("got record response: %v, want record response: %v", record.Response, wantRecordThree.Response)
		}

	})

}

var testConfig = config.Records{
	Headers: config.HeadersToRecord{
		Exclude: []string{"Header-To-Exclude"},
	},
}

var payload = RecordPayload{
	Hash:       "testHash",
	Curl:       "testCurl",
	Scheme:     "https",
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

var wantRecordThree = Record{
	Request: Request{
		Hash:       "testHash",
		Curl:       "testCurl",
		Scheme:     "https",
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
	Response: Response{
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
