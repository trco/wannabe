package actions

import (
	"reflect"
	"testing"
	"wannabe/utils"

	"github.com/clbanning/mxj"
)

func TestEncodeBody(t *testing.T) {
	tests := []struct {
		name            string
		decodedBody     interface{}
		contentType     []string
		contentEncoding []string
		want            []byte
		wantErr         bool
	}{
		{
			name:            "json content type",
			decodedBody:     map[string]string{"key": "value"},
			contentType:     []string{"application/json"},
			contentEncoding: []string{},
			want:            []byte(`{"key":"value"}`),
			wantErr:         false,
		},
		{
			name:            "xml content type",
			decodedBody:     map[string]interface{}{"root": map[string]interface{}{"key": "value"}},
			contentType:     []string{"application/xml"},
			contentEncoding: []string{},
			want:            nil,
			wantErr:         false,
		},
		{
			name:            "text/plain content type",
			decodedBody:     "plain text",
			contentType:     []string{"text/plain"},
			contentEncoding: []string{},
			want:            []byte("plain text"),
			wantErr:         false,
		},
		{
			name:            "text/html content type",
			decodedBody:     "\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e",
			contentType:     []string{"text/html"},
			contentEncoding: []string{},
			want:            []byte("\u003c!DOCTYPE html\u003e\n \u003chtml\u003e\n \u003ch1\u003eHerman Melville - Moby-Dick\u003c/h1\u003e\n \u003c/html\u003e"),
			wantErr:         false,
		},
		{
			name:            "unsupported content type",
			decodedBody:     "unsupported",
			contentType:     []string{"unsupported"},
			contentEncoding: []string{},
			want:            nil,
			wantErr:         true,
		},
		{
			name:            "json content type, gzip content encoding",
			decodedBody:     map[string]string{"key": "value"},
			contentType:     []string{"application/json"},
			contentEncoding: []string{"gzip"},
			want:            getCompressedEncodedBody([]byte(`{"key":"value"}`)),
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeBody(tt.decodedBody, tt.contentEncoding, tt.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "xml content type" {
				mapValue := mxj.Map(tt.decodedBody.(map[string]interface{}))
				body, _ := mapValue.Xml()

				tt.want = body
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCompressedEncodedBody(data []byte) []byte {
	compressedEncodedBody, _ := utils.Gzip(data)
	return compressedEncodedBody
}
