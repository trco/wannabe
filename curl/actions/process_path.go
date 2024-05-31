package actions

import (
	"fmt"
	"strings"
	"wannabe/curl/utils"
	"wannabe/types"
)

func ProcessPath(path string, config types.Path) (string, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	if path == "" {
		return "", nil
	}

	pathParts := strings.Split(strippedPath, "/")

	utils.SetWildcardsByIndex(pathParts, config.Wildcards)
	rebuiltPath := "/" + strings.Join(pathParts, "/")

	processedPath, err := utils.ReplaceRegexPatterns(rebuiltPath, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessPath: failed compiling regex: %v", err)
	}

	return processedPath, nil
}
