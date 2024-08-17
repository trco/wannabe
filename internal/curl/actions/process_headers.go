package actions

import (
	"github.com/trco/wannabe/curl/utils"
	"github.com/trco/wannabe/types"
)

func ProcessHeaders(headersMap map[string][]string, config types.Headers) []types.Header {
	headers := utils.FilterHeadersToInclude(headersMap, config.Include)

	wildcards := config.Wildcards
	for _, wildcard := range wildcards {
		if !utils.KeyExists(headers, wildcard.Key) {
			// TODO log warning
			continue
		}

		utils.SetPlaceholderByKey(headers, wildcard)
	}

	headerSlice := utils.HeadersMapToSlice(headers)
	sortedHeaders := utils.SortHeaderSlice(headerSlice)

	return sortedHeaders
}
