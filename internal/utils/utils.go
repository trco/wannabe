package utils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func GetContentType(contentTypeHeader []string) string {
	switch {
	case sliceItemContains(contentTypeHeader, "application/json"):
		return "application/json"
	case sliceItemContains(contentTypeHeader, "application/xml"):
		return "application/xml"
	case sliceItemContains(contentTypeHeader, "text/xml"):
		return "text/xml"
	case sliceItemContains(contentTypeHeader, "text/plain"):
		return "text/plain"
	case sliceItemContains(contentTypeHeader, "text/html"):
		return "text/html"
	default:
		return ""
	}
}

func GetContentEncoding(contentEncodingHeader []string) string {
	switch {
	case sliceItemContains(contentEncodingHeader, "gzip"):
		return "gzip"
	default:
		return ""
	}
}

// Gunzip decompresses a gzip-compressed byte slice.
func Gunzip(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Gunzip: failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	decompressedData, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("Gunzip: failed to read decompressed data: %v", err)
	}

	return decompressedData, nil
}

// Gzip compresses a byte slice using gzip.
func Gzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	_, err := gzWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("Gzip: failed to write data to gzip writer: %v", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("Gzip: failed to close gzip writer: %v", err)
	}

	return buf.Bytes(), nil
}

func sliceItemContains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}
