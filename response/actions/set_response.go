package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy/proxyutil"
	"github.com/clbanning/mxj"
)

func SetResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	var record types.Record

	err := json.Unmarshal(encodedRecord, &record)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed unmarshaling record: %v", err)
	}

	statusCode := record.Response.StatusCode

	decodedBody := record.Response.Body
	contentTypeResponse := record.Response.Headers["Content-Type"]
	body, err := prepareBody(decodedBody, contentTypeResponse)

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

func sliceItemContains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}

func prepareBody(decodedBody interface{}, contentType []string) ([]byte, error) {
	if sliceItemContains(contentType, "application/json") {
		body, err := json.Marshal(decodedBody)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed marshaling body: %v", err)
		}

		return body, nil
	} else if sliceItemContains(contentType, "application/xml") || sliceItemContains(contentType, "text/xml") {
		mapValue := mxj.Map(decodedBody.(map[string]interface{}))
		body, err := mapValue.Xml()
		if err != nil {
			return body, fmt.Errorf("GenerateRecord: failed unmarshaling XML body: %v", err)
		}

		return body, nil
	}

	// else if sliceItemContains(contentType, "text/plain") {
	// 	body = string(encodedBody)

	// 	return body, nil
	// } else {
	// 	return nil, fmt.Errorf("SetResponse: unsupported content type: %s", contentType)
	// }

	return nil, nil
}
