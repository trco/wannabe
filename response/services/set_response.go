package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"wannabe/response/actions"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

func SetResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	var record types.Record

	err := json.Unmarshal(encodedRecord, &record)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed unmarshaling record: %v", err)
	}

	decodedBody := record.Response.Body
	contentTypeHeader := record.Response.Headers["Content-Type"]
	contentEncodingHeader := record.Response.Headers["Content-Encoding"]
	body, err := actions.EncodeBody(decodedBody, contentTypeHeader, contentEncodingHeader)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed encoding response body: %v", err)
	}

	statusCode := record.Response.StatusCode

	response := proxyutil.NewResponse(statusCode, bytes.NewReader(body), request)

	headers := record.Response.Headers
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
