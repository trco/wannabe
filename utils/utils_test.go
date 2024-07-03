package utils

import (
	"bytes"
	"compress/gzip"
	"reflect"
	"testing"
)

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
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(originalData); err != nil {
		t.Fatalf("Failed to write gzip data: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("Failed to close gzip writer: %v", err)
	}
	compressedData := buf.Bytes()

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
