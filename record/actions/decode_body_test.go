package actions

import (
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
)

type TestCaseDecodeBody struct {
	encodedBody []byte
	contentType []string
	expected    interface{}
}

func TestDecodeBody(t *testing.T) {
	testCases := map[string]TestCaseDecodeBody{
		"empty body": {
			encodedBody: []byte{},
			contentType: []string{"application/json"},
			expected:    nil,
		},
		"empty content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{},
			expected:    nil,
		},
		"json content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"application/json"},
			expected:    map[string]interface{}{"key": "value"},
		},
		"xml content type": {
			encodedBody: []byte(`<root><key>value</key></root>`),
			contentType: []string{"application/xml"},
			expected:    nil,
		},
		"plain text content type": {
			encodedBody: []byte("plain text"),
			contentType: []string{"text/plain"},
			expected:    "plain text",
		},
		"unsupported content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"unsupported/type"},
			expected:    nil,
		},
	}

	for testKey, tc := range testCases {
		decodedBody, _ := DecodeBody(tc.encodedBody, tc.contentType)

		if testKey == "xml content type" {
			xmlMap, _ := mxj.NewMapXml(tc.encodedBody)
			tc.expected = xmlMap
		}

		if !reflect.DeepEqual(tc.expected, decodedBody) {
			t.Errorf("failed test case: %v, expected decoded body: %v, actual decoded body: %v", testKey, tc.expected, decodedBody)
		}
	}
}
