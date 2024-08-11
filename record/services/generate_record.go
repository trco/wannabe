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

	requestContentEncoding := payload.RequestHeaders[contentEncoding]
	requestContentType := payload.RequestHeaders[contentType]
	requestBody, err := actions.DecodeBody(payload.RequestBody, requestContentEncoding, requestContentType)
	if err != nil {
		return nil, err
	}

	responseContentEncoding := payload.ResponseHeaders[contentEncoding]
	responseContentType := payload.ResponseHeaders[contentType]
	responseBody, err := actions.DecodeBody(payload.ResponseBody, responseContentEncoding, responseContentType)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	record := types.Record{
		Request: types.Request{
			Hash:       payload.Hash,
			Curl:       payload.Curl,
			Scheme:     payload.Scheme,
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
