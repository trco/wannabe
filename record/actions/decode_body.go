package actions

import (
	"encoding/json"
	"fmt"

	"wannabe/utils"

	"github.com/clbanning/mxj"
)

func DecodeBody(encodedBody []byte, contentTypeHeader []string) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	contentType := utils.GetContentType(contentTypeHeader)

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
	case contentType == "text/plain", contentType == "text/html":
		body = string(encodedBody)
	default:
		return body, fmt.Errorf("DecodeBody: unsupported content type: %s", contentType)
	}

	return body, nil
}
