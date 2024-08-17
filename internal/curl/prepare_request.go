package curl

import (
	"bytes"
	"fmt"
	"net/http"
)

func PrepareRequest(httpMethod string, url string, body string) (*http.Request, error) {
	var request *http.Request
	var err error

	if body == "" {
		request, err = http.NewRequest(httpMethod, url, nil)
	} else {
		bodyBuffer := bytes.NewBufferString(body)
		request, err = http.NewRequest(httpMethod, url, bodyBuffer)
	}

	if err != nil {
		return nil, fmt.Errorf("PrepareRequest: failed generating request: %v", err)
	}

	return request, nil
}
