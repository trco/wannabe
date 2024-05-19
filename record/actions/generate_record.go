package actions

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/record/entities"
	"wannabe/types"
)

func GenerateRecord(config types.Records, payload entities.RecordPayload) ([]byte, error) {
	requestHeaders := filterRequestHeaders(payload.RequestHeaders, config.Headers.Exclude)

	requestBody, err := prepareBody(payload.RequestBody)
	if err != nil {
		return nil, err
	}

	// FIXME case when response body is text should also be handled
	responseBody, err := prepareBody(payload.ResponseBody)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	record := entities.Record{
		Request: entities.Request{
			Hash:       payload.Hash,
			Curl:       payload.Curl,
			HttpMethod: payload.HttpMethod,
			Host:       payload.Host,
			Path:       payload.Path,
			Query:      payload.Query,
			Headers:    requestHeaders,
			Body:       requestBody,
		},
		Response: entities.Response{
			StatusCode: payload.StatusCode,
			Headers:    payload.ResponseHeaders,
			Body:       responseBody,
		},
		Metadata: entities.Metadata{
			GeneratedAt: entities.Timestamp{
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

func prepareBody(encodedBody []byte) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	err := json.Unmarshal(encodedBody, &body)
	if err != nil {
		return body, fmt.Errorf("GenerateRecord: failed unmarshaling body: %v", err)
	}

	return body, nil
}
