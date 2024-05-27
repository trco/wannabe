package actions

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
	"wannabe/types"
)

// general

func setWildcardsByIndex(slice []string, wildcards []types.WildcardIndex) {
	for _, wildcard := range wildcards {
		if isIndexOutOfBounds(slice, *wildcard.Index) {
			// TODO log warning
			continue
		}

		setPlaceholderByIndex(slice, wildcard)
	}
}

func setWildcardsByKey(inputMap map[string]string, wildcards []types.WildcardKey) {
	for _, wildcard := range wildcards {
		if !keyExists(inputMap, wildcard.Key) {
			// TODO log warning
			continue
		}

		setPlaceholderByKey(inputMap, wildcard)
	}
}

func setPlaceholderByIndex(parts []string, wildcard types.WildcardIndex) {
	if wildcard.Placeholder != "" {
		parts[*wildcard.Index] = wildcard.Placeholder
	} else {
		parts[*wildcard.Index] = "{wannabe}"
	}
}

func setPlaceholderByKey(inputMap map[string]string, wildcard types.WildcardKey) {
	if wildcard.Placeholder != "" {
		inputMap[wildcard.Key] = wildcard.Placeholder
	} else {
		inputMap[wildcard.Key] = "{wannabe}"
	}
}

func replaceRegexPatterns(processedString string, regexes []types.Regex, isQuery bool) (string, error) {
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
			regex.Placeholder = "{wannabe}"
		}

		if isQuery {
			regex.Placeholder = url.QueryEscape(regex.Placeholder)
		}

		processedString = compiledPattern.ReplaceAllString(processedString, regex.Placeholder)
	}

	return processedString, nil
}

// query

func mapValuesToSingleString(queryMap map[string][]string) map[string]string {
	query := make(map[string]string)
	for key := range queryMap {
		queryValue := strings.Join(queryMap[key], ",")
		query[key] = queryValue
	}

	return query
}

func buildQuery(query map[string]string) string {
	values := url.Values{}
	for k, v := range query {
		values.Add(k, v)
	}

	// values.Encode sorts params alphabetically by key
	return "?" + values.Encode()
}

// headers

func filterHeadersToInclude(headersMap map[string][]string, headersToInclude []string) map[string]string {
	if len(headersToInclude) == 0 {
		for key := range headersMap {
			headersToInclude = append(headersToInclude, key)
		}
	}

	headers := make(map[string]string)
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

func headersMapToSlice(headersMap map[string]string) []types.Header {
	var headerSlice []types.Header
	for key, value := range headersMap {
		headerSlice = append(headerSlice, types.Header{Key: key, Value: value})
	}

	return headerSlice
}

func sortHeaderSlice(headerSlice []types.Header) []types.Header {
	sort.Slice(headerSlice, func(i, j int) bool {
		return headerSlice[i].Key < headerSlice[j].Key
	})

	return headerSlice
}

// checks

func isIndexOutOfBounds[T interface{}](slice []T, index int) bool {
	return index < 0 || index >= len(slice)
}

func keyExists[T interface{}](stringsMap map[string]T, key string) bool {
	_, exists := stringsMap[key]
	return exists
}
