package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/config"
)

func ProcessBody(body []byte, config config.Body) (string, error) {
	if len(body) == 0 {
		return "", nil
	}

	var bodyMap map[string]interface{}

	// results in alphabetically ordered json
	json.Unmarshal(body, &bodyMap)

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return "", fmt.Errorf("ProcessBody: failed marshaling request body: %v", err)
	}

	bodyString := string(bodyBytes)

	processedBodyString, err := replaceRegexPatterns(bodyString, config.Regexes)
	if err != nil {
		return "", fmt.Errorf("ProcessBody: failed compiling regex: %v", err)
	}

	return processedBodyString, nil
}
