package actions

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"wannabe/types"
)

func GenerateCurlPayload(request *http.Request) (types.CurlPayload, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return types.CurlPayload{}, fmt.Errorf("PrepareGenerateCurlPayload: failed reading request body: %v", err)
	}
	defer request.Body.Close()

	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	curlPayload := types.CurlPayload{
		HttpMethod:     request.Method,
		Host:           request.URL.Host,
		Path:           request.URL.Path,
		Query:          request.URL.Query(),
		RequestHeaders: request.Header,
		RequestBody:    body,
	}

	return curlPayload, nil
}
