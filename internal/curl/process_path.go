package curl

import (
	"fmt"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func ProcessPath(path string, config config.Path) (string, error) {
	strippedPath := strings.TrimPrefix(path, "/")
	if path == "" {
		return "", nil
	}

	pathParts := strings.Split(strippedPath, "/")

	SetWildcardsByIndex(pathParts, config.Wildcards)
	rebuiltPath := "/" + strings.Join(pathParts, "/")

	processedPath, err := ReplaceRegexPatterns(rebuiltPath, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessPath: failed compiling regex: %v", err)
	}

	return processedPath, nil
}
