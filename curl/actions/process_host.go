package actions

import (
	"fmt"
	"strings"
	"wannabe/config"
)

func ProcessHost(host string, config config.Host) (string, error) {
	strippedHost := strings.TrimPrefix(strings.TrimPrefix(host, "http://"), "https://")
	hostParts := strings.Split(strippedHost, ".")

	setWildcardsByIndex(hostParts, config.Wildcards)
	rebuiltHost := strings.Join(hostParts, ".")

	processedHost, err := replaceRegexPatterns(rebuiltHost, config.Regexes)
	if err != nil {
		return "", fmt.Errorf("ProcessHost: failed compiling regex: %v", err)
	}

	return processedHost, nil
}
