package utils

import (
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
			got := GetContentType(tt.header)
			if got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
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
			got := sliceItemContains(tt.slice, tt.value)
			if got != tt.want {
				t.Errorf("sliceItemContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
