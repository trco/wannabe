package record

import (
	"bytes"
	"io"

	"github.com/trco/wannabe/internal/handlers/utils"
)

type RecordPayload struct {
	Hash            string
	Curl            string
	Scheme          string
	HttpMethod      string
	Host            string
	Path            string
	Query           map[string][]string
	RequestHeaders  map[string][]string
	RequestBody     []byte
	StatusCode      int
	ResponseHeaders map[string][]string
	ResponseBody    []byte
}

func GenerateRecordPayload(session utils.Session, hash string, curl string) (RecordPayload, error) {
	request := session.Request()
	response := session.Response()

	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return RecordPayload{}, err
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return RecordPayload{}, err
	}
	defer response.Body.Close()

	// set body back to response
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	recordPayload := RecordPayload{
		Hash:            hash,
		Curl:            curl,
		Scheme:          request.URL.Scheme,
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
