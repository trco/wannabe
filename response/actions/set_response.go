package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"wannabe/record/entities"

	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

func SetResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	var record entities.Record

	err := json.Unmarshal(encodedRecord, &record)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed unmarshaling record: %v", err)
	}

	statusCode := record.Response.StatusCode
	headers := record.Response.Headers

	body, err := json.Marshal(record.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed marshaling body: %v", err)
	}

	response := proxyutil.NewResponse(statusCode, bytes.NewReader(body), request)
	setHeaders(response, headers)

	return response, nil
}

func setHeaders(response *http.Response, headers map[string][]string) {
	for key, value := range headers {
		for _, v := range value {
			response.Header.Set(key, v)
		}
	}
}
