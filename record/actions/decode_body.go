package actions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/mxj"
)

func DecodeBody(encodedBody []byte, contentType []string) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	// REVIEW is this reasonable? what if Content-Type header is not present? enforce Content-Type header and validate for its presence?
	if len(contentType) == 0 {
		return body, nil
	}

	if sliceItemContains(contentType, "application/json") {
		err := json.Unmarshal(encodedBody, &body)
		if err != nil {
			return body, fmt.Errorf("DecodeBody: failed unmarshaling JSON body: %v", err)
		}
	} else if sliceItemContains(contentType, "application/xml") || sliceItemContains(contentType, "text/xml") {
		xmlMap, err := mxj.NewMapXml(encodedBody)
		if err != nil {
			return body, fmt.Errorf("DecodeBody: failed unmarshaling XML body: %v", err)
		}
		body = xmlMap
	} else if sliceItemContains(contentType, "text/plain") {
		body = string(encodedBody)
	} else {
		return body, fmt.Errorf("DecodeBody: unsupported content type: %s", contentType)
	}

	return body, nil
}

func sliceItemContains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}
