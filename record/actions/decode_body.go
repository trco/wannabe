package actions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/mxj"
)

func DecodeBody(encodedBody []byte, contentTypeHeader []string) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	contentType := getContentType(contentTypeHeader)

	switch {
	case contentType == "application/json":
		err := json.Unmarshal(encodedBody, &body)
		if err != nil {
			return body, fmt.Errorf("DecodeBody: failed unmarshaling JSON body: %v", err)
		}
	case contentType == "application/xml", contentType == "text/xml":
		xmlMap, err := mxj.NewMapXml(encodedBody)
		if err != nil {
			return body, fmt.Errorf("DecodeBody: failed unmarshaling XML body: %v", err)
		}
		body = xmlMap
	case contentType == "text/plain":
		body = string(encodedBody)
	default:
		return body, fmt.Errorf("DecodeBody: unsupported content type: %s", contentType)
	}

	return body, nil
}

func getContentType(contentTypeHeader []string) string {
	switch {
	case sliceItemContains(contentTypeHeader, "application/json"):
		return "application/json"
	case sliceItemContains(contentTypeHeader, "application/xml"):
		return "application/xml"
	case sliceItemContains(contentTypeHeader, "text/xml"):
		return "text/xml"
	case sliceItemContains(contentTypeHeader, "text/plain"):
		return "text/plain"
	default:
		return ""
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
