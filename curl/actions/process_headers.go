package actions

import "wannabe/types"

func ProcessHeaders(headersMap map[string][]string, config types.Headers) []Header {
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
