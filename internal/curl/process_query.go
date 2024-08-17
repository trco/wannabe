package curl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func ProcessQuery(queryMap map[string][]string, config config.Query) (string, error) {
	if len(queryMap) == 0 {
		return "", nil
	}

	query := MapValuesToSingleString(queryMap)

	SetWildcardsByKey(query, config.Wildcards)
	rebuiltQuery := BuildQuery(query)

	processedQuery, err := ReplaceRegexPatterns(rebuiltQuery, config.Regexes, true)
	if err != nil {
		return "", fmt.Errorf("ProcessQuery: failed compiling regex: %v", err)
	}

	return processedQuery, nil
}

func SetWildcardsByKey(inputMap map[string]string, wildcards []config.WildcardKey) {
	for _, wildcard := range wildcards {
		if !KeyExists(inputMap, wildcard.Key) {
			// TODO log warning
			continue
		}

		SetPlaceholderByKey(inputMap, wildcard)
	}
}

func KeyExists[T interface{}](stringsMap map[string]T, key string) bool {
	_, exists := stringsMap[key]
	return exists
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
