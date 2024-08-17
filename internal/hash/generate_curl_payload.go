package hash

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type curlPayload struct {
	scheme         string
	httpMethod     string
	host           string
	path           string
	query          map[string][]string
	requestHeaders map[string][]string
	requestBody    []byte
}

func generateCurlPayload(request *http.Request) (curlPayload, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return curlPayload{}, fmt.Errorf("PrepareGenerateCurlPayload: failed reading request body: %v", err)
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	curlPayload := curlPayload{
		scheme:         request.URL.Scheme,
		httpMethod:     request.Method,
		host:           request.URL.Host,
		path:           request.URL.Path,
		query:          request.URL.Query(),
		requestHeaders: request.Header,
		requestBody:    body,
	}

	return curlPayload, nil
}
