package services

import (
	"reflect"
	"testing"
	"time"
	"wannabe/record/entities"
)

func TestExtractRecords(t *testing.T) {
	records, _ := ExtractRecords(requestBody)
	record := records[0]

	if !reflect.DeepEqual(expectedRecordB, record) {
		t.Errorf("expected record: %v, actual record: %v", expectedRecordB, record)
	}
}

// reusable variables
var requestBody = []byte{91, 10, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 114, 101, 113, 117, 101, 115, 116, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 104, 116, 116, 112, 77, 101, 116, 104, 111, 100, 34, 58, 32, 34, 80, 79, 83, 84, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 104, 111, 115, 116, 34, 58, 32, 34, 104, 116, 116, 112, 115, 58, 47, 47, 97, 110, 97, 108, 121, 116, 105, 99, 115, 100, 97, 116, 97, 46, 103, 111, 111, 103, 108, 101, 97, 112, 105, 115, 46, 99, 111, 109, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 112, 97, 116, 104, 34, 58, 32, 34, 116, 101, 115, 116, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 113, 117, 101, 114, 121, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 104, 101, 97, 100, 101, 114, 115, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 115, 111, 110, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 98, 111, 100, 121, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 32, 32, 32, 32, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 34, 114, 101, 115, 112, 111, 110, 115, 101, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 115, 116, 97, 116, 117, 115, 67, 111, 100, 101, 34, 58, 32, 50, 48, 48, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 104, 101, 97, 100, 101, 114, 115, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 106, 115, 111, 110, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 93, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 125, 44, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 98, 111, 100, 121, 34, 58, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 116, 101, 115, 116, 34, 58, 32, 34, 116, 101, 115, 116, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 125, 10, 93}

var expectedRecordB = entities.Record{
	Request: entities.Request{
		HttpMethod: "POST",
		Host:       "https://analyticsdata.googleapis.com",
		Path:       "test",
		Query: map[string]string{
			"test": "test",
		},
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: map[string]interface{}{
			"test": "test",
		},
	},
	Response: entities.Response{
		StatusCode: 200,
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: map[string]interface{}{
			"test": "test",
		},
	},
	Metadata: entities.Metadata{
		RequestedAt: entities.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
		GeneratedAt: entities.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
		RegeneratedAt: entities.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
	},
}
