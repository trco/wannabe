package utils

import "strings"

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
	default:
		return ""
	}
}

func sliceItemContains(slice []string, value string) bool {
	for _, item := range slice {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}
