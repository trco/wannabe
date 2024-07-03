package actions

import (
	"reflect"
	"testing"
	"wannabe/utils"

	"github.com/clbanning/mxj"
)

func TestDecodeBody(t *testing.T) {
	tests := []struct {
		name            string
		encodedBody     []byte
		contentType     []string
		contentEncoding []string
		want            interface{}
	}{
		{
			name:            "empty body",
			encodedBody:     []byte{},
			contentType:     []string{"application/json"},
			contentEncoding: []string{},
			want:            nil,
		},
		{
			name:            "empty content type",
			encodedBody:     []byte(`{"key": "value"}`),
			contentType:     []string{},
			contentEncoding: []string{},
			want:            nil,
		},
		{
			name:            "json content type",
			encodedBody:     []byte(`{"key": "value"}`),
			contentType:     []string{"application/json"},
			contentEncoding: []string{},
			want:            map[string]interface{}{"key": "value"},
		},
		{
			name:            "xml content type",
			encodedBody:     []byte(`<root><key>value</key></root>`),
			contentType:     []string{"application/xml"},
			contentEncoding: []string{},
			want:            nil,
		},
		{
			name:            "plain text content type",
			encodedBody:     []byte("plain text"),
			contentType:     []string{"text/plain"},
			contentEncoding: []string{},
			want:            "plain text",
		},
		{
			name:            "html text content type",
			encodedBody:     []byte("\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e"),
			contentType:     []string{"text/html"},
			contentEncoding: []string{},
			want:            "\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e",
		},
		{
			name:            "unsupported content type",
			encodedBody:     []byte(`{"key": "value"}`),
			contentType:     []string{"unsupported/type"},
			contentEncoding: []string{},
			want:            nil,
		},
		{
			name:            "json content type, gzip content encoding",
			encodedBody:     getCompressedEncodedBody([]byte(`{"key": "value"}`)),
			contentType:     []string{"application/json"},
			contentEncoding: []string{"gzip"},
			want:            map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := DecodeBody(tt.encodedBody, tt.contentType, tt.contentEncoding)

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

func getCompressedEncodedBody(data []byte) []byte {
	compressedEncodedBody, _ := utils.Gzip(data)
	return compressedEncodedBody
}
