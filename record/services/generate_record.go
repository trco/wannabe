package services

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/record/actions"
	"wannabe/types"
)

var contentType = "Content-Type"
var contentEncoding = "Content-Encoding"

func GenerateRecord(config types.Records, payload types.RecordPayload) ([]byte, error) {
	requestHeaders := actions.FilterHeaders(payload.RequestHeaders, config.Headers.Exclude)

	requestContentType := payload.RequestHeaders[contentType]
	requestContentEncoding := payload.RequestHeaders[contentEncoding]
	requestBody, err := actions.DecodeBody(payload.RequestBody, requestContentType, requestContentEncoding)
	if err != nil {
		return nil, err
	}

	responseContentType := payload.ResponseHeaders[contentType]
	responseContentEncoding := payload.ResponseHeaders[contentEncoding]
	responseBody, err := actions.DecodeBody(payload.ResponseBody, responseContentType, responseContentEncoding)
	if err != nil {
		return nil, err
	}

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
			Body:       requestBody,
		},
		Response: types.Response{
			StatusCode: payload.StatusCode,
			Headers:    payload.ResponseHeaders,
			Body:       responseBody,
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
