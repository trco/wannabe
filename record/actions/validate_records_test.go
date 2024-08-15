package actions

import (
	"reflect"
	"testing"
	"time"

	"github.com/trco/wannabe/types"
)

func TestValidateRecords(t *testing.T) {
	tests := []struct {
		name    string
		record  types.Record
		wantErr []string
	}{
		{
			name:    "valid record",
			record:  validRecord,
			wantErr: []string{""},
		},
		{
			name:    "invalid record",
			record:  invalidRecord,
			wantErr: []string{"Key: 'Record.Request.HttpMethod' Error:Field validation for 'HttpMethod' failed on the 'oneof' tag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, _ := ValidateRecords([]types.Record{tt.record})

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("ValidateRecords() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

var validRecord = types.Record{
	Request: types.Request{
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
	Metadata: types.Metadata{
		GeneratedAt: types.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
		RegeneratedAt: types.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
	},
}

var invalidRecord = types.Record{
	Request: types.Request{
		Hash:       "testHash",
		Curl:       "testCurl",
		Scheme:     "https",
		HttpMethod: "INVALID",
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
	Metadata: types.Metadata{
		GeneratedAt: types.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
		RegeneratedAt: types.Timestamp{
			Unix: 0,
			UTC:  time.Time{},
		},
	},
}
