package actions

import (
	"fmt"
	"strings"
	"wannabe/config"
)

func ProcessPath(path string, config config.Path) (string, error) {
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
