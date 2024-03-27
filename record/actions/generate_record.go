package actions

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/config"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
)

// generates record from ctx *fiber.Ctx request and response, server, hash and curl
func GenerateRecord(ctx *fiber.Ctx, config config.Records, server string, curl string, hash string) ([]byte, error) {
	requestHeaders := filterRequestHeaders(ctx.GetReqHeaders(), config.Headers.Exclude)

	requestBody, err := prepareBody(ctx.Body())
	if err != nil {
		return nil, err
	}

	responseBody, err := prepareBody(ctx.Response().Body())
	if err != nil {
		return nil, err
	}

	record := entities.Record{
		Request: entities.Request{
			Hash:       hash,
			Curl:       curl,
			HttpMethod: ctx.Method(),
			Host:       server,
			Path:       ctx.Path(),
			Query:      ctx.Queries(),
			Headers:    requestHeaders,
			Body:       requestBody,
		},
		Response: entities.Response{
			StatusCode: ctx.Response().StatusCode(),
			Headers:    ctx.GetRespHeaders(),
			Body:       responseBody,
		},
		Metadata: entities.Metadata{
			RequestedAt: entities.Timestamp{
				Unix: ctx.Context().Time().Unix(),
				UTC:  ctx.Context().Time().UTC(),
			},
			GeneratedAt: entities.Timestamp{
				Unix: time.Now().Unix(),
				UTC:  time.Now().UTC(),
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

	err := json.Unmarshal(encodedBody, &body)
	if err != nil {
		return body, fmt.Errorf("GenerateRecord: failed unmarshaling body: %v", err)
	}

	return body, nil
}
