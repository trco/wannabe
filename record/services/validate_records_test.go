package services

import (
	"testing"
	"time"
	"wannabe/config"
	"wannabe/record/entities"
)

func TestValidateRecords(t *testing.T) {
	// valid record
	validationErrors, _ := ValidateRecords(testConfigB, []entities.Record{validRecord})

	if validationErrors[0] != "" {
		t.Errorf("validation failed although it should not, validationErrors: %v", validationErrors)
	}

	// invalid record
	validationErrors, _ = ValidateRecords(testConfigB, []entities.Record{invalidRecord})
	expectedErrs := "Key: 'Record.Request.Host' Error:Field validation for 'Host' failed on the 'host_not_matching_config_server' tag"

	if validationErrors[0] != expectedErrs {
		t.Errorf("validation succeeded although it should not, validationErrors: %v", validationErrors)
	}
}

// reusable variables

var testConfigB = config.Config{
	Mode:            "server",
	FailOnReadError: true,
	StorageProvider: config.StorageProvider{
		Type:       "filesystem",
		Regenerate: false,
		FilesystemConfig: config.FilesystemConfig{
			Folder:           "records",
			RegenerateFolder: "",
			Format:           "json",
		},
	},
	Server: "https://analyticsdata.googleapis.com",
}

var validRecord = entities.Record{
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

var invalidRecord = entities.Record{
	Request: entities.Request{
		HttpMethod: "POST",
		Host:       "https://test.com",
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
