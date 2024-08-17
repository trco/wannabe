package curl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func processQuery(queryMap map[string][]string, config config.Query) (string, error) {
	if len(queryMap) == 0 {
		return "", nil
	}

	query := mapValuesToSingleString(queryMap)

	setWildcardsByKey(query, config.Wildcards)
	rebuiltQuery := buildQuery(query)

	processedQuery, err := replaceRegexPatterns(rebuiltQuery, config.Regexes, true)
	if err != nil {
		return "", fmt.Errorf("processQuery: failed compiling regex: %v", err)
	}

	return processedQuery, nil
}

func setWildcardsByKey(inputMap map[string]string, wildcards []config.WildcardKey) {
	for _, wildcard := range wildcards {
		if !keyExists(inputMap, wildcard.Key) {
			// TODO log warning
			continue
		}

		setPlaceholderByKey(inputMap, wildcard)
	}
}

func keyExists[T interface{}](stringsMap map[string]T, key string) bool {
	_, exists := stringsMap[key]
	return exists
}

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
