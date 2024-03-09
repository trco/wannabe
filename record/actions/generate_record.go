package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/record/entities"

	"github.com/gofiber/fiber/v2"
)

// generates record from ctx *fiber.Ctx request and response, server, hash and curl
func GenerateRecord(ctx *fiber.Ctx, server string, curl string, hash string) ([]byte, error) {

	requestBody, err := prepareBody(ctx.Body())
	if err != nil {
		return nil, err
	}

	responseBody, err := prepareBody(ctx.Response().Body())
	if err != nil {
		return nil, err
	}

	record := entities.Record{
		// TODO add Metadata (timestamp,...)
		Request: entities.Request{
			Hash:       hash,
			HttpMethod: ctx.Method(),
			Host:       server,
			Path:       ctx.Path(),
			Query:      ctx.Queries(),
			// Query:      string(ctx.Request().URI().QueryString()),
			Headers: ctx.GetReqHeaders(),
			Body:    requestBody,
			Curl:    curl,
		},
		Response: entities.Response{
			StatusCode: ctx.Response().StatusCode(),
			Headers:    ctx.GetRespHeaders(),
			Body:       responseBody,
		},
	}

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("GenerateRecord: failed marshaling record: %v", err)
	}

	return recordBytes, nil
}

func prepareBody(bodyBytes []byte) (entities.BodyMap, error) {
	var body entities.BodyMap

	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return body, fmt.Errorf("GenerateRecord: failed unmarshaling body: %v", err)
	}

	return body, nil
}
