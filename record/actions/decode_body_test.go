package actions

import (
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
)

func TestDecodeBody(t *testing.T) {
	tests := []struct {
		name        string
		encodedBody []byte
		contentType []string
		want        interface{}
	}{
		{
			name:        "empty body",
			encodedBody: []byte{},
			contentType: []string{"application/json"},
			want:        nil,
		},
		{
			name:        "empty content type",
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{},
			want:        nil,
		},
		{
			name:        "json content type",
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"application/json"},
			want:        map[string]interface{}{"key": "value"},
		},
		{
			name:        "xml content type",
			encodedBody: []byte(`<root><key>value</key></root>`),
			contentType: []string{"application/xml"},
			want:        nil,
		},
		{
			name:        "plain text content type",
			encodedBody: []byte("plain text"),
			contentType: []string{"text/plain"},
			want:        "plain text",
		},
		{
			name:        "unsupported content type",
			encodedBody: []byte(`{"key": "value"}`),
			contentType: []string{"unsupported/type"},
			want:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := DecodeBody(tt.encodedBody, tt.contentType)

			if tt.name == "xml content type" {
				xmlMap, _ := mxj.NewMapXml(tt.encodedBody)
				tt.want = xmlMap
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
