package actions

import (
	"testing"
	"time"
	"wannabe/types"
)

func TestValidateRecords(t *testing.T) {
	// valid record
	validationErrors, _ := ValidateRecords([]types.Record{record})

	if validationErrors[0] != "" {
		t.Errorf("valid record validated as invalid")
	}

	record.Request.HttpMethod = "INVALID"

	// invalid record
	validationErrors, _ = ValidateRecords([]types.Record{record})

	if validationErrors[0] == "" {
		t.Errorf("invalid record validated as valid")
	}
}

var record = types.Record{
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
