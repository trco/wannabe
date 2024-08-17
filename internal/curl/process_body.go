package curl

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/internal/config"
)

func processBody(requestBody []byte, config config.Body) (string, error) {
	if len(requestBody) == 0 || string(requestBody) == "null" {
		return "", nil
	}

	var body interface{}

	// results in alphabetically ordered json
	json.Unmarshal(requestBody, &body)

	requestBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("processBody: failed marshaling request body: %v", err)
	}

	bodyString := string(requestBody)

	processedBodyString, err := replaceRegexPatterns(bodyString, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("processBody: failed compiling regex: %v", err)
	}

	return processedBodyString, nil
}
