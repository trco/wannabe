package actions

import (
	"bytes"
	"fmt"
	"net/http"
)

func PrepareRequest(httpMethod string, url string, body string) (*http.Request, error) {
	bodyBuffer := bytes.NewBufferString(body)

	request, err := http.NewRequest(httpMethod, url, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("PrepareRequest: failed generating request: %v", err)
	}

	return request, nil
}
