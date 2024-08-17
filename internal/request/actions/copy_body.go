package actions

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func CopyBody(request *http.Request) (io.ReadCloser, error) {
	var requestBody io.ReadCloser

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("CopyBody: failed reading request body: %v", err)
	}
	defer request.Body.Close()

	requestBody = io.NopCloser(bytes.NewReader(body))
	// set body back to request
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	return requestBody, nil
}
