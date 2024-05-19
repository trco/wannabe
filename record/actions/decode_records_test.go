package actions

import (
	"reflect"
	"testing"
	"wannabe/record/entities"
)

func TestDecodeRecords(t *testing.T) {
	decodedRecords, _ := DecodeRecords(encodedRecords)
	record := decodedRecords[0]

	if !reflect.DeepEqual(expectedRecordA.Request, record.Request) {
		t.Errorf("expected record request: %v, actual record request: %v", expectedRecordA.Request, record.Request)
	}

	if !reflect.DeepEqual(expectedRecordA.Response, record.Response) {
		t.Errorf("expected record response: %v, actual record response: %v", expectedRecordA.Response, record.Response)
	}
}

// reusable variables
var expectedRecordA = entities.Record{
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

var encodedRecords = [][]byte{
	{123, 34, 114, 101, 113, 117, 101, 115, 116, 34, 58, 123, 34, 104, 97, 115, 104, 34, 58, 34, 116, 101, 115, 116, 72, 97, 115, 104,
		34, 44, 34, 99, 117, 114, 108, 34, 58, 34, 116, 101, 115, 116, 67, 117, 114, 108, 34, 44, 34, 104, 116, 116, 112, 77, 101, 116,
		104, 111, 100, 34, 58, 34, 80, 79, 83, 84, 34, 44, 34, 104, 111, 115, 116, 34, 58, 34, 116, 101, 115, 116, 46, 99, 111, 109, 34,
		44, 34, 112, 97, 116, 104, 34, 58, 34, 47, 116, 101, 115, 116, 34, 44, 34, 113, 117, 101, 114, 121, 34, 58, 123, 34, 116, 101,
		115, 116, 34, 58, 91, 34, 116, 101, 115, 116, 34, 93, 125, 44, 34, 104, 101, 97, 100, 101, 114, 115, 34, 58, 123, 34, 65, 99, 99,
		101, 112, 116, 34, 58, 91, 34, 116, 101, 115, 116, 34, 93, 44, 34, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 34, 58,
		91, 34, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 115, 111, 110, 34, 93, 125, 44, 34, 98, 111, 100, 121, 34, 58,
		123, 34, 116, 101, 115, 116, 34, 58, 34, 116, 101, 115, 116, 34, 125, 125, 44, 34, 114, 101, 115, 112, 111, 110, 115, 101, 34, 58,
		123, 34, 115, 116, 97, 116, 117, 115, 67, 111, 100, 101, 34, 58, 50, 48, 48, 44, 34, 104, 101, 97, 100, 101, 114, 115, 34, 58, 123,
		34, 65, 99, 99, 101, 112, 116, 34, 58, 91, 34, 116, 101, 115, 116, 34, 93, 44, 34, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112,
		101, 34, 58, 91, 34, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 115, 111, 110, 34, 93, 125, 44, 34, 98, 111, 100, 121,
		34, 58, 123, 34, 116, 101, 115, 116, 34, 58, 34, 116, 101, 115, 116, 34, 125, 125, 44, 34, 109, 101, 116, 97, 100, 97, 116, 97, 34, 58,
		123, 34, 103, 101, 110, 101, 114, 97, 116, 101, 100, 65, 116, 34, 58, 123, 34, 117, 110, 105, 120, 34, 58, 48, 44, 34, 117, 116, 99, 34,
		58, 34, 48, 48, 48, 49, 45, 48, 49, 45, 48, 49, 84, 48, 48, 58, 48, 48, 58, 48, 48, 90, 34, 125, 44, 34, 114, 101, 103, 101, 110, 101,
		114, 97, 116, 101, 100, 65, 116, 34, 58, 123, 34, 117, 110, 105, 120, 34, 58, 48, 44, 34, 117, 116, 99, 34, 58, 34, 48, 48, 48, 49, 45,
		48, 49, 45, 48, 49, 84, 48, 48, 58, 48, 48, 58, 48, 48, 90, 34, 125, 125, 125},
}
