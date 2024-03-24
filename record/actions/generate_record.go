package actions

import (
	"encoding/json"
	"fmt"
	"time"
	"wannabe/config"
	"wannabe/record/common"

	"github.com/gofiber/fiber/v2"
)

// generates record from ctx *fiber.Ctx request and response, server, hash and curl
func GenerateRecord(ctx *fiber.Ctx, config config.Records, server string, curl string, hash string) ([]byte, error) {
	requestHeaders := common.FilterRequestHeaders(ctx.GetReqHeaders(), config.Headers.Exclude)

	requestBody, err := common.PrepareBody(ctx.Body())
	if err != nil {
		return nil, fmt.Errorf("GenerateRecord: failed unmarshaling request body: %v", err)
	}

	responseBody, err := common.PrepareBody(ctx.Response().Body())
	if err != nil {
		return nil, fmt.Errorf("GenerateRecord: failed unmarshaling response body: %v", err)
	}

	record := common.Record{
		Request: common.Request{
			Hash:       hash,
			Curl:       curl,
			HttpMethod: ctx.Method(),
			Host:       server,
			Path:       ctx.Path(),
			Query:      ctx.Queries(),
			Headers:    requestHeaders,
			Body:       requestBody,
		},
		Response: common.Response{
			StatusCode: ctx.Response().StatusCode(),
			Headers:    ctx.GetRespHeaders(),
			Body:       responseBody,
		},
		Metadata: common.Metadata{
			RequestedAt: common.Timestamp{
				Unix: ctx.Context().Time().Unix(),
				UTC:  ctx.Context().Time().UTC(),
			},
			GeneratedAt: common.Timestamp{
				Unix: time.Now().Unix(),
				UTC:  time.Now().UTC(),
			},
		},
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("GenerateRecord: failed marshaling record: %v", err)
	}

	return recordBytes, nil
}
