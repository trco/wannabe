package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/config"
)

func ProcessBody(bodyBytes []byte, config config.Body) (string, error) {
	if len(bodyBytes) == 0 {
		return "", nil
	}

	var body interface{}

	// results in alphabetically ordered json
	json.Unmarshal(bodyBytes, &body)

	bodyBytes, err := json.Marshal(body)
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
