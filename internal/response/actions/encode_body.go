package actions

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/utils"

	"github.com/clbanning/mxj"
)

func EncodeBody(decodedBody interface{}, contentTypeHeader []string, contentEncodingHeader []string) ([]byte, error) {
	var body []byte

	contentType := utils.GetContentType(contentTypeHeader)
	switch {
	case contentType == "application/json":
		var err error
		body, err = json.Marshal(decodedBody)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed marshaling body: %v", err)
		}
	case contentType == "application/xml", contentType == "text/xml":
		var err error
		mapValue := mxj.Map(decodedBody.(map[string]interface{}))
		body, err = mapValue.Xml()
		if err != nil {
			return body, fmt.Errorf("GenerateRecord: failed unmarshaling XML body: %v", err)
		}
	case contentType == "text/plain", contentType == "text/html":
		body = []byte(decodedBody.(string))
	default:
		return nil, fmt.Errorf("SetResponse: unsupported content type: %s", contentTypeHeader)
	}

	// gzip the body that was unzipped before storing it to the record
	contentEncoding := utils.GetContentEncoding(contentEncodingHeader)
	if contentEncoding == "gzip" {
		compressedBody, err := utils.Gzip(body)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed compressing response body: %s", err)
		}

		return compressedBody, nil
	}

	return body, nil
}
