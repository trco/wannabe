package actions

import (
	"fmt"
	"io"
	"net/http"
	"wannabe/curl/entities"
)

func GenerateCurlPayload(originalRequest *http.Request) (entities.CurlPayload, error) {
	body, err := io.ReadAll(originalRequest.Body)
	if err != nil {
		return entities.CurlPayload{}, fmt.Errorf("PrepareGenerateCurlPayload: failed reading request body: %v", err)
	}
	defer originalRequest.Body.Close()

	curlPayload := entities.CurlPayload{
		HttpMethod:     originalRequest.Method,
		Host:           originalRequest.URL.Host,
		Path:           originalRequest.URL.Path,
		Query:          originalRequest.URL.Query(),
		RequestHeaders: originalRequest.Header,
		RequestBody:    body,
	}

	return curlPayload, nil
}
