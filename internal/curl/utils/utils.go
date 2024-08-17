package utils

import (
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/trco/wannabe/types"
)

func SetWildcardsByIndex(slice []string, wildcards []types.WildcardIndex) {
	for _, wildcard := range wildcards {
		if IsIndexOutOfBounds(slice, *wildcard.Index) {
			// TODO log warning
			continue
		}

		SetPlaceholderByIndex(slice, wildcard)
	}
}

func SetWildcardsByKey(inputMap map[string]string, wildcards []types.WildcardKey) {
	for _, wildcard := range wildcards {
		if !KeyExists(inputMap, wildcard.Key) {
			// TODO log warning
			continue
		}

		SetPlaceholderByKey(inputMap, wildcard)
	}
}

func SetPlaceholderByIndex(parts []string, wildcard types.WildcardIndex) {
	if wildcard.Placeholder != "" {
		parts[*wildcard.Index] = wildcard.Placeholder
	} else {
		parts[*wildcard.Index] = "wannabe"
	}
}

func SetPlaceholderByKey(inputMap map[string]string, wildcard types.WildcardKey) {
	if wildcard.Placeholder != "" {
		inputMap[wildcard.Key] = wildcard.Placeholder
	} else {
		inputMap[wildcard.Key] = "wannabe"
	}
}

func ReplaceRegexPatterns(processedString string, regexes []types.Regex, isQuery bool) (string, error) {
	for _, regex := range regexes {
		compiledPattern, err := regexp.Compile(regex.Pattern)
		if err != nil {
			return "", err
		}

		match := compiledPattern.MatchString(processedString)
		if !match {
			// TODO log warning
			continue
		}

		if regex.Placeholder == "" {
			regex.Placeholder = "wannabe"
		}

		if isQuery {
			regex.Placeholder = url.QueryEscape(regex.Placeholder)
		}

		processedString = compiledPattern.ReplaceAllString(processedString, regex.Placeholder)
	}

	return processedString, nil
}

func MapValuesToSingleString(queryMap map[string][]string) map[string]string {
	query := make(map[string]string)
	for key := range queryMap {
		queryValue := strings.Join(queryMap[key], ",")
		query[key] = queryValue
	}

	return query
}

func BuildQuery(query map[string]string) string {
	values := url.Values{}
	for k, v := range query {
		values.Add(k, v)
	}

	// values.Encode sorts params alphabetically by key
	return "?" + values.Encode()
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

func HeadersMapToSlice(headersMap map[string]string) []types.Header {
	var headerSlice []types.Header
	for key, value := range headersMap {
		headerSlice = append(headerSlice, types.Header{Key: key, Value: value})
	}

	return headerSlice
}

func SortHeaderSlice(headerSlice []types.Header) []types.Header {
	sort.Slice(headerSlice, func(i, j int) bool {
		return headerSlice[i].Key < headerSlice[j].Key
	})

	return headerSlice
}

func IsIndexOutOfBounds[T interface{}](slice []T, index int) bool {
	return index < 0 || index >= len(slice)
}

func KeyExists[T interface{}](stringsMap map[string]T, key string) bool {
	_, exists := stringsMap[key]
	return exists
}
