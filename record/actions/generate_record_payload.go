package actions

import (
	"bytes"
	"io"
	"wannabe/types"
)

func GenerateRecordPayload(wannabeSession types.WannabeSession, hash string, curl string) (types.RecordPayload, error) {
	request := wannabeSession.Request()
	response := wannabeSession.Response()

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return types.RecordPayload{}, err
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return types.RecordPayload{}, err
	}
	defer response.Body.Close()

	// set body back to response
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	recordPayload := types.RecordPayload{
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
