package record

import (
	"reflect"
	"testing"

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
			got, _ := DecodeBody(tt.encodedBody, tt.contentEncoding, tt.contentType)

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
	compressedEncodedBody, _ := Gzip(data)
	return compressedEncodedBody
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name   string
		header []string
		want   string
	}{
		{
			name:   "json content type",
			header: []string{"application/json"},
			want:   "application/json",
		},
		{
			name:   "xml content type",
			header: []string{"application/xml"},
			want:   "application/xml",
		},
		{
			name:   "text/xml content type",
			header: []string{"text/xml"},
			want:   "text/xml",
		},
		{
			name:   "text/plain content type",
			header: []string{"text/plain"},
			want:   "text/plain",
		},
		{
			name:   "text/html content type",
			header: []string{"text/html"},
			want:   "text/html",
		},
		{
			name:   "unknown content type",
			header: []string{"unknown"},
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContentType(tt.header); got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContentEncoding(t *testing.T) {
	tests := []struct {
		name   string
		header []string
		want   string
	}{
		{
			name:   "gzip encoding",
			header: []string{"gzip"},
			want:   "gzip",
		},
		{
			name:   "empty encoding",
			header: []string{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContentEncoding(tt.header); got != tt.want {
				t.Errorf("GetContentEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGunzip(t *testing.T) {
	originalData := []byte("test data")
	compressedData, _ := Gzip(originalData)

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr string
	}{
		{
			name:    "valid gzip data",
			data:    compressedData,
			want:    originalData,
			wantErr: "",
		},
		{
			name:    "invalid gzip data",
			data:    []byte("invalid"),
			want:    nil,
			wantErr: "Gunzip: failed to create gzip reader: unexpected EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gunzip(tt.data)
			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("Gunzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gunzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGzip(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{
			name:    "valid data",
			data:    []byte("test data"),
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gzip(tt.data)
			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("Gzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Decompress to verify the content
			decompressedData, err := Gunzip(got)
			if err != nil {
				t.Errorf("Gunzip() failed to decompress: %v", err)
				return
			}
			if !reflect.DeepEqual(decompressedData, tt.data) {
				t.Errorf("Gzip() decompressed data = %v, want %v", decompressedData, tt.data)
			}
		})
	}
}

func TestSliceItemContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		value string
		want  bool
	}{
		{
			name:  "slice contains value",
			slice: []string{"apple", "banana", "cherry"},
			value: "banana",
			want:  true,
		},
		{
			name:  "slice does not contain value",
			slice: []string{"apple", "banana", "cherry"},
			value: "grape",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sliceItemContains(tt.slice, tt.value); got != tt.want {
				t.Errorf("sliceItemContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
