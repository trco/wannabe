package actions

import (
	"fmt"
	"strings"
	"wannabe/types"
)

func ProcessHost(host string, config types.Host) (string, error) {
	hostParts := strings.Split(host, ".")

	setWildcardsByIndex(hostParts, config.Wildcards)
	rebuiltHost := strings.Join(hostParts, ".")

	processedHost, err := replaceRegexPatterns(rebuiltHost, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessHost: failed compiling regex: %v", err)
	}

	return processedHost, nil
}
