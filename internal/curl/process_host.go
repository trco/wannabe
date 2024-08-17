package curl

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/trco/wannabe/internal/config"
)

func ProcessHost(host string, config config.Host) (string, error) {
	hostParts := strings.Split(host, ".")

	SetWildcardsByIndex(hostParts, config.Wildcards)
	rebuiltHost := strings.Join(hostParts, ".")

	processedHost, err := ReplaceRegexPatterns(rebuiltHost, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessHost: failed compiling regex: %v", err)
	}

	return processedHost, nil
}

func SetWildcardsByIndex(slice []string, wildcards []config.WildcardIndex) {
	for _, wildcard := range wildcards {
		if IsIndexOutOfBounds(slice, *wildcard.Index) {
			// TODO log warning
			continue
		}

		SetPlaceholderByIndex(slice, wildcard)
	}
}

func IsIndexOutOfBounds[T interface{}](slice []T, index int) bool {
	return index < 0 || index >= len(slice)
}

func SetPlaceholderByIndex(parts []string, wildcard config.WildcardIndex) {
	if wildcard.Placeholder != "" {
		parts[*wildcard.Index] = wildcard.Placeholder
	} else {
		parts[*wildcard.Index] = "wannabe"
	}
}

func ReplaceRegexPatterns(processedString string, regexes []config.Regex, isQuery bool) (string, error) {
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
