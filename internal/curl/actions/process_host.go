package actions

import (
	"fmt"
	"strings"

	"github.com/trco/wannabe/curl/utils"
	"github.com/trco/wannabe/types"
)

func ProcessHost(host string, config types.Host) (string, error) {
	hostParts := strings.Split(host, ".")

	utils.SetWildcardsByIndex(hostParts, config.Wildcards)
	rebuiltHost := strings.Join(hostParts, ".")

	processedHost, err := utils.ReplaceRegexPatterns(rebuiltHost, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessHost: failed compiling regex: %v", err)
	}

	return processedHost, nil
}
