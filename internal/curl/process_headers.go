package curl

import (
	"sort"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

type Header struct {
	Key   string
	Value string
}

func ProcessHeaders(headersMap map[string][]string, config config.Headers) []Header {
	headers := FilterHeadersToInclude(headersMap, config.Include)

	wildcards := config.Wildcards
	for _, wildcard := range wildcards {
		if !KeyExists(headers, wildcard.Key) {
			// TODO log warning
			continue
		}

		SetPlaceholderByKey(headers, wildcard)
	}

	headerSlice := HeadersMapToSlice(headers)
	sortedHeaders := SortHeaderSlice(headerSlice)

	return sortedHeaders
}

func FilterHeadersToInclude(headersMap map[string][]string, headersToInclude []string) map[string]string {
	headers := make(map[string]string)

	if len(headersToInclude) == 0 {
		return headers
	}

	for _, key := range headersToInclude {
		if !KeyExists(headersMap, key) {
			// TODO log warning
			continue
		}

		headerValue := strings.Join(headersMap[key], ",")
		headers[key] = headerValue
	}

	return headers
}

func SetPlaceholderByKey(inputMap map[string]string, wildcard config.WildcardKey) {
	if wildcard.Placeholder != "" {
		inputMap[wildcard.Key] = wildcard.Placeholder
	} else {
		inputMap[wildcard.Key] = "wannabe"
	}
}

func HeadersMapToSlice(headersMap map[string]string) []Header {
	var headerSlice []Header
	for key, value := range headersMap {
		headerSlice = append(headerSlice, Header{Key: key, Value: value})
	}

	return headerSlice
}

func SortHeaderSlice(headerSlice []Header) []Header {
	sort.Slice(headerSlice, func(i, j int) bool {
		return headerSlice[i].Key < headerSlice[j].Key
	})

	return headerSlice
}
