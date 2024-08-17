package hash

import (
	"sort"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

type header struct {
	key   string
	value string
}

func processHeaders(headersMap map[string][]string, config config.Headers) []header {
	headers := filterHeadersToInclude(headersMap, config.Include)

	wildcards := config.Wildcards
	for _, wildcard := range wildcards {
		if !keyExists(headers, wildcard.Key) {
			// TODO log warning
			continue
		}

		setPlaceholderByKey(headers, wildcard)
	}

	headerSlice := headersMapToSlice(headers)
	sortedHeaders := sortHeaderSlice(headerSlice)

	return sortedHeaders
}

func filterHeadersToInclude(headersMap map[string][]string, headersToInclude []string) map[string]string {
	headers := make(map[string]string)

	if len(headersToInclude) == 0 {
		return headers
	}

	for _, key := range headersToInclude {
		if !keyExists(headersMap, key) {
			// TODO log warning
			continue
		}

		headerValue := strings.Join(headersMap[key], ",")
		headers[key] = headerValue
	}

	return headers
}

func setPlaceholderByKey(inputMap map[string]string, wildcard config.WildcardKey) {
	if wildcard.Placeholder != "" {
		inputMap[wildcard.Key] = wildcard.Placeholder
	} else {
		inputMap[wildcard.Key] = "wannabe"
	}
}

func headersMapToSlice(headersMap map[string]string) []header {
	var headerSlice []header
	for key, value := range headersMap {
		headerSlice = append(headerSlice, header{key: key, value: value})
	}

	return headerSlice
}

func sortHeaderSlice(headerSlice []header) []header {
	sort.Slice(headerSlice, func(i, j int) bool {
		return headerSlice[i].key < headerSlice[j].key
	})

	return headerSlice
}
