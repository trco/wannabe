package actions

import "wannabe/config"

func ProcessHeaders(headersMap map[string][]string, config config.Headers) []Header {
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
