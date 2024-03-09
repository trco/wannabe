package actions

import (
	"fmt"
	"strings"
	"wannabe/config"
)

func ProcessPath(path string, config config.Path) (string, error) {
	if path == "" {
		return path, nil
	}

	strippedPath := strings.TrimPrefix(path, "/")
	pathParts := strings.Split(strippedPath, "/")

	setWildcardsByIndex(pathParts, config.Wildcards)
	rebuiltPath := "/" + strings.Join(pathParts, "/")

	processedPath, err := replaceRegexPatterns(rebuiltPath, config.Regexes)
	if err != nil {
		return "", fmt.Errorf("ProcessPath: failed compiling regex: %v", err)
	}

	return processedPath, nil
}
