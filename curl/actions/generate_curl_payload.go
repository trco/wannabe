package actions

import (
	"fmt"
	"io"
	"net/http"
	"wannabe/types"
)

func GenerateCurlPayload(originalRequest *http.Request) (types.CurlPayload, error) {
	body, err := io.ReadAll(originalRequest.Body)
	if err != nil {
		return types.CurlPayload{}, fmt.Errorf("PrepareGenerateCurlPayload: failed reading request body: %v", err)
	}
	defer originalRequest.Body.Close()

	curlPayload := types.CurlPayload{
		HttpMethod:     originalRequest.Method,
		Host:           originalRequest.URL.Host,
		Path:           originalRequest.URL.Path,
		Query:          originalRequest.URL.Query(),
		RequestHeaders: originalRequest.Header,
		RequestBody:    body,
	}

	return curlPayload, nil
}
