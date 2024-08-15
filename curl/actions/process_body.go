package actions

import (
	"encoding/json"
	"fmt"

	"github.com/trco/wannabe/curl/utils"
	"github.com/trco/wannabe/types"
)

func ProcessBody(requestBody []byte, config types.Body) (string, error) {
	if len(requestBody) == 0 || string(requestBody) == "null" {
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

	processedBodyString, err := utils.ReplaceRegexPatterns(bodyString, config.Regexes, false)
	if err != nil {
		return "", fmt.Errorf("ProcessBody: failed compiling regex: %v", err)
	}

	return processedBodyString, nil
}
