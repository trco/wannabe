package curl

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func processHost(host string, config config.Host) (string, error) {
	hostParts := strings.Split(host, ".")

	setWildcardsByIndex(hostParts, config.Wildcards)
	rebuiltHost := strings.Join(hostParts, ".")

	processedHost, err := replaceRegexPatterns(rebuiltHost, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("processHost: failed compiling regex: %v", err)
	}

	return processedHost, nil
}

func setWildcardsByIndex(slice []string, wildcards []config.WildcardIndex) {
	for _, wildcard := range wildcards {
		if isIndexOutOfBounds(slice, *wildcard.Index) {
			// TODO log warning
			continue
		}

		setPlaceholderByIndex(slice, wildcard)
	}
}

func isIndexOutOfBounds[T interface{}](slice []T, index int) bool {
	return index < 0 || index >= len(slice)
}

func setPlaceholderByIndex(parts []string, wildcard config.WildcardIndex) {
	if wildcard.Placeholder != "" {
		parts[*wildcard.Index] = wildcard.Placeholder
	} else {
		parts[*wildcard.Index] = "wannabe"
	}
}

func replaceRegexPatterns(processedString string, regexes []config.Regex, isQuery bool) (string, error) {
	for _, regex := range regexes {
		compiledPattern, err := regexp.Compile(regex.Pattern)
		if err != nil {
			return "", err
		}

		match := compiledPattern.MatchString(processedString)
		if !match {
			// TODO log warning
			continue
		}

		if regex.Placeholder == "" {
			regex.Placeholder = "wannabe"
		}

		if isQuery {
			regex.Placeholder = url.QueryEscape(regex.Placeholder)
		}

		processedString = compiledPattern.ReplaceAllString(processedString, regex.Placeholder)
	}

	return processedString, nil
}
