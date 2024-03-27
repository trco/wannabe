package actions

import (
	"encoding/json"
	"fmt"
	"wannabe/config"
)

func ProcessBody(requestBody []byte, config config.Body) (string, error) {
	if len(requestBody) == 0 {
		return "", nil
	}

	var body interface{}

	// results in alphabetically ordered json
	json.Unmarshal(requestBody, &body)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("ProcessBody: failed marshaling request body: %v", err)
	}

	bodyString := string(requestBody)

	processedBodyString, err := replaceRegexPatterns(bodyString, config.Regexes)
	if err != nil {
		return "", fmt.Errorf("ProcessBody: failed compiling regex: %v", err)
	}

	return processedBodyString, nil
}
