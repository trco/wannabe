package actions

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/trco/wannabe/types"
)

func GenerateRequest(recordRequest types.Request) (*http.Request, error) {
	if recordRequest.Path == "" {
		recordRequest.Path = "/"
	}

	body := recordRequest.Body
	var requestBody []byte

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(recordRequest.HttpMethod, recordRequest.Path, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	request.URL.Scheme = recordRequest.Scheme
	request.URL.Host = recordRequest.Host

	query := request.URL.Query()
	for key, value := range recordRequest.Query {
		for _, item := range value {
			query.Add(key, item)
		}
	}
	request.URL.RawQuery = query.Encode()

	for key, value := range recordRequest.Headers {
		for _, item := range value {
			request.Header.Set(key, item)
		}
	}

	return request, nil
}
