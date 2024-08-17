package hash

import (
	"fmt"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func processPath(path string, config config.Path) (string, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	if path == "" {
		return "", nil
	}

	pathParts := strings.Split(strippedPath, "/")

	setWildcardsByIndex(pathParts, config.Wildcards)
	rebuiltPath := "/" + strings.Join(pathParts, "/")

	processedPath, err := replaceRegexPatterns(rebuiltPath, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("processPath: failed compiling regex: %v", err)
	}

	return processedPath, nil
}
