package actions

import (
	"reflect"
	"testing"

	"github.com/clbanning/mxj"
)

func TestEncodeBody(t *testing.T) {
	tests := []struct {
		name        string
		decodedBody interface{}
		contentType []string
		want        []byte
		wantErr     bool
	}{
		{
			name:        "json content type",
			decodedBody: map[string]string{"key": "value"},
			contentType: []string{"application/json"},
			want:        []byte(`{"key":"value"}`),
			wantErr:     false,
		},
		{
			name:        "xml content type",
			decodedBody: map[string]interface{}{"root": map[string]interface{}{"key": "value"}},
			contentType: []string{"application/xml"},
			want:        nil,
			wantErr:     false,
		},
		{
			name:        "text/plain content type",
			decodedBody: "plain text",
			contentType: []string{"text/plain"},
			want:        []byte("plain text"),
			wantErr:     false,
		},
		{
			name:        "unsupported content type",
			decodedBody: "unsupported",
			contentType: []string{"unsupported"},
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeBody(tt.decodedBody, tt.contentType)
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
