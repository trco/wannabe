package actions

import (
	"fmt"
	"wannabe/config"
)

func ProcessQuery(query map[string]string, config config.Query) (string, error) {
	if len(query) == 0 {
		return "", nil
	}

	setWildcardsByKey(query, config.Wildcards)
	rebuiltQuery := buildQuery(query)

	processedQuery, err := replaceRegexPatterns(rebuiltQuery, config.Regexes)
	if err != nil {
		return "", fmt.Errorf("ProcessQuery: failed compiling regex: %v", err)
	}

	return processedQuery, nil
}
