package actions

import (
	"fmt"
	"wannabe/curl/utils"
	"wannabe/types"
)

func ProcessQuery(queryMap map[string][]string, config types.Query) (string, error) {
	if len(queryMap) == 0 {
		return "", nil
	}

	query := utils.MapValuesToSingleString(queryMap)

	utils.SetWildcardsByKey(query, config.Wildcards)
	rebuiltQuery := utils.BuildQuery(query)

	processedQuery, err := utils.ReplaceRegexPatterns(rebuiltQuery, config.Regexes, true)
	if err != nil {
		return "", fmt.Errorf("ProcessQuery: failed compiling regex: %v", err)
	}

	return processedQuery, nil
}
