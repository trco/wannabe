package actions

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/types"
)

func GenerateRecord(config types.Records, payload types.RecordPayload) ([]byte, error) {
	requestHeaders := filterRequestHeaders(payload.RequestHeaders, config.Headers.Exclude)
	timestamp := time.Now()

	record := types.Record{
		Request: types.Request{
			Hash:       payload.Hash,
			Curl:       payload.Curl,
			HttpMethod: payload.HttpMethod,
			Host:       payload.Host,
			Path:       payload.Path,
			Query:      payload.Query,
			Headers:    requestHeaders,
			Body:       payload.RequestBody,
		},
		Response: types.Response{
			StatusCode: payload.StatusCode,
			Headers:    payload.ResponseHeaders,
			Body:       payload.ResponseBody,
		},
		Metadata: types.Metadata{
			GeneratedAt: types.Timestamp{
				Unix: timestamp.Unix(),
				UTC:  timestamp.UTC(),
			},
		},
	}

	encodedRecord, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("GenerateRecord: failed marshaling record: %v", err)
	}

	return encodedRecord, nil
}

func filterRequestHeaders(headers map[string][]string, exclude []string) map[string][]string {
	filteredRequestHeaders := make(map[string][]string)

	for key, values := range headers {
		if !contains(exclude, key) {
			filteredRequestHeaders[key] = values
		}
	}

	return filteredRequestHeaders
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
