package actions

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wannabe/types"

	"github.com/clbanning/mxj"
)

func GenerateRecord(config types.Records, payload types.RecordPayload) ([]byte, error) {
	requestHeaders := filterRequestHeaders(payload.RequestHeaders, config.Headers.Exclude)

	requestContentType := payload.RequestHeaders["Content-Type"]
	requestBody, err := prepareBody(payload.RequestBody, requestContentType)
	if err != nil {
		return nil, err
	}

	responseContentType := payload.ResponseHeaders["Content-Type"]
	responseBody, err := prepareBody(payload.ResponseBody, responseContentType)
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

func sliceItemContains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}

func prepareBody(encodedBody []byte, contentType []string) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	// REVIEW is this reasonable? what if Content-Type header is not present? enforce Content-Type header and validate for its presence?
	if len(contentType) == 0 {
		return body, nil
	}

	if sliceItemContains(contentType, "application/json") {
		err := json.Unmarshal(encodedBody, &body)
		if err != nil {
			return body, fmt.Errorf("GenerateRecord: failed unmarshaling JSON body: %v", err)
		}
	} else if sliceItemContains(contentType, "application/xml") || sliceItemContains(contentType, "text/xml") {
		xmlMap, err := mxj.NewMapXml(encodedBody)
		if err != nil {
			return body, fmt.Errorf("GenerateRecord: failed unmarshaling XML body: %v", err)
		}
		body = xmlMap
	} else if sliceItemContains(contentType, "text/plain") {
		body = string(encodedBody)
	} else {
		return body, fmt.Errorf("GenerateRecord: unsupported content type: %s", contentType)
	}

	return body, nil
}
