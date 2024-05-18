package actions

import (
	"fmt"
	"wannabe/config"
)

func ProcessQuery(queryMap map[string][]string, config config.Query) (string, error) {
	if len(queryMap) == 0 {
		return "", nil
	}

	query := mapValuesToSingleString(queryMap)

	setWildcardsByKey(query, config.Wildcards)
	rebuiltQuery := buildQuery(query)

	processedQuery, err := replaceRegexPatterns(rebuiltQuery, config.Regexes, true)
	if err != nil {
		return "", fmt.Errorf("ProcessQuery: failed compiling regex: %v", err)
	}

	return processedQuery, nil
}
