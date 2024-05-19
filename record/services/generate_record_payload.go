package services

import (
	"bytes"
	"io"
	"wannabe/record/entities"
	"wannabe/types"
)

func GenerateRecordPayload(wannabeSession types.WannabeSession, hash string, curl string) (entities.RecordPayload, error) {
	request := wannabeSession.Request()
	response := wannabeSession.Response()

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return entities.RecordPayload{}, err
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return entities.RecordPayload{}, err
	}
	defer response.Body.Close()

	// set body back to response
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	recordPayload := entities.RecordPayload{
		Hash:            hash,
		Curl:            curl,
		HttpMethod:      request.Method,
		Host:            request.URL.Host,
		Path:            request.URL.Path,
		Query:           request.URL.Query(),
		RequestHeaders:  request.Header,
		RequestBody:     requestBody,
		StatusCode:      response.StatusCode,
		ResponseHeaders: response.Header,
		ResponseBody:    responseBody,
	}

	return recordPayload, nil
}
