package actions

import (
	"encoding/json"
	"fmt"

	"wannabe/utils"

	"github.com/clbanning/mxj"
)

func EncodeBody(decodedBody interface{}, contentTypeHeader []string) ([]byte, error) {
	contentType := utils.GetContentType(contentTypeHeader)

	switch {
	case contentType == "application/json":
		body, err := json.Marshal(decodedBody)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed marshaling body: %v", err)
		}

		return body, nil
	case contentType == "application/xml", contentType == "text/xml":
		mapValue := mxj.Map(decodedBody.(map[string]interface{}))
		body, err := mapValue.Xml()
		if err != nil {
			return body, fmt.Errorf("GenerateRecord: failed unmarshaling XML body: %v", err)
		}

		return body, nil
	case contentType == "text/plain":
		body := []byte(decodedBody.(string))

		return body, nil
	default:
		return nil, fmt.Errorf("SetResponse: unsupported content type: %s", contentTypeHeader)
	}
}
