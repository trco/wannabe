package actions

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/utils"

	"github.com/clbanning/mxj"
)

func DecodeBody(encodedBody []byte, contentEncodingHeader []string, contentTypeHeader []string) (interface{}, error) {
	var body interface{}

	if len(encodedBody) == 0 {
		return body, nil
	}

	contentEncoding := utils.GetContentEncoding(contentEncodingHeader)

	// gunzip the body before decoding it to the interface
	if contentEncoding == "gzip" {
		var err error
		encodedBody, err = utils.Gunzip(encodedBody)
		if err != nil {
			return nil, fmt.Errorf("DecodeBody: failed unzipping body: %v", err)
		}
	}

	contentType := utils.GetContentType(contentTypeHeader)

	switch {
	case contentType == "application/json":
		err := json.Unmarshal(encodedBody, &body)
		if err != nil {
			return nil, fmt.Errorf("DecodeBody: failed unmarshaling JSON body: %v", err)
		}
	case contentType == "application/xml", contentType == "text/xml":
		xmlMap, err := mxj.NewMapXml(encodedBody)
		if err != nil {
			return nil, fmt.Errorf("DecodeBody: failed unmarshaling XML body: %v", err)
		}
		body = xmlMap
	case contentType == "text/plain", contentType == "text/html":
		body = string(encodedBody)
	default:
		return nil, fmt.Errorf("DecodeBody: unsupported content type: %s", contentType)
	}

	return body, nil
}
