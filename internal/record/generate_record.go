package record

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/trco/wannabe/internal/config"
)

var contentType = "Content-Type"
var contentEncoding = "Content-Encoding"

type Record struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Metadata Metadata `json:"metadata"`
}

type Request struct {
	Hash       string              `json:"hash"`
	Curl       string              `json:"curl"`
	Scheme     string              `json:"scheme" validate:"required,oneof=http https"`
	HttpMethod string              `json:"httpMethod" validate:"required,oneof=GET POST PUT DELETE PATCH HEAD CONNECT OPTIONS TRACE"`
	Host       string              `json:"host" validate:"required"`
	Path       string              `json:"path"`
	Query      map[string][]string `json:"query"`
	Headers    map[string][]string `json:"headers"`
	Body       interface{}         `json:"body" validate:"required_if=HttpMethod POST,required_if=HttpMethod PUT,required_if=HttpMethod PATCH"`
}

type Response struct {
	StatusCode int                 `json:"statusCode" validate:"required"`
	Headers    map[string][]string `json:"headers" validate:"content_type_header_present"`
	Body       interface{}         `json:"body" validate:"required"`
}

type Metadata struct {
	GeneratedAt   Timestamp `json:"generatedAt"`
	RegeneratedAt Timestamp `json:"regeneratedAt"`
}

type Timestamp struct {
	Unix int64     `json:"unix"`
	UTC  time.Time `json:"utc"`
}

func GenerateRecord(config config.Records, payload RecordPayload) ([]byte, error) {
	requestHeaders := FilterHeaders(payload.RequestHeaders, config.Headers.Exclude)

	requestContentEncoding := payload.RequestHeaders[contentEncoding]
	requestContentType := payload.RequestHeaders[contentType]
	requestBody, err := DecodeBody(payload.RequestBody, requestContentEncoding, requestContentType)
	if err != nil {
		return nil, err
	}

	responseContentEncoding := payload.ResponseHeaders[contentEncoding]
	responseContentType := payload.ResponseHeaders[contentType]
	responseBody, err := DecodeBody(payload.ResponseBody, responseContentEncoding, responseContentType)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	record := Record{
		Request: Request{
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
		Response: Response{
			StatusCode: payload.StatusCode,
			Headers:    payload.ResponseHeaders,
			Body:       responseBody,
		},
		Metadata: Metadata{
			GeneratedAt: Timestamp{
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
