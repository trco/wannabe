package actions

import (
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
)

func TestDecodeBody(t *testing.T) {
	testCases := map[string]struct {
		encodedBody []byte
		contentType []string
		want        interface{}
	}{
		"empty body": {
			encodedBody: []byte{},
			contentType: []string{"application/json"},
			want:        nil,
		},
		"empty content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{},
			want:        nil,
		},
		"json content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"application/json"},
			want:        map[string]interface{}{"key": "value"},
		},
		"xml content type": {
			encodedBody: []byte(`<root><key>value</key></root>`),
			contentType: []string{"application/xml"},
			want:        nil,
		},
		"plain text content type": {
			encodedBody: []byte("plain text"),
			contentType: []string{"text/plain"},
			want:        "plain text",
		},
		"unsupported content type": {
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"unsupported/type"},
			want:        nil,
		},
	}

	for testKey, tc := range testCases {
		t.Run(testKey, func(t *testing.T) {
			got, _ := DecodeBody(tc.encodedBody, tc.contentType)

			if testKey == "xml content type" {
				xmlMap, _ := mxj.NewMapXml(tc.encodedBody)
				tc.want = xmlMap
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("DecodeBody() = %v, want %v", got, tc.want)
			}
		})
	}
}
