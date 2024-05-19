package actions

import (
	"fmt"
	"strings"
	"wannabe/types"
)

func ProcessPath(path string, config types.Path) (string, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	if path == "" {
		return "", nil
	}

	pathParts := strings.Split(strippedPath, "/")

	setWildcardsByIndex(pathParts, config.Wildcards)
	rebuiltPath := "/" + strings.Join(pathParts, "/")

	processedPath, err := replaceRegexPatterns(rebuiltPath, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessPath: failed compiling regex: %v", err)
	}

	return processedPath, nil
}
