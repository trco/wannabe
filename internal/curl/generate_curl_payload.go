package curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type CurlPayload struct {
	Scheme         string
	HttpMethod     string
	Host           string
	Path           string
	Query          map[string][]string
	RequestHeaders map[string][]string
	RequestBody    []byte
}

func GenerateCurlPayload(request *http.Request) (CurlPayload, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return CurlPayload{}, fmt.Errorf("PrepareGenerateCurlPayload: failed reading request body: %v", err)
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	curlPayload := CurlPayload{
		Scheme:         request.URL.Scheme,
		HttpMethod:     request.Method,
		Host:           request.URL.Host,
		Path:           request.URL.Path,
		Query:          request.URL.Query(),
		RequestHeaders: request.Header,
		RequestBody:    body,
	}

	return curlPayload, nil
}
