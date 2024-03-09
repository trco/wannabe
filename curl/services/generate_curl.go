package services

import (
	"wannabe/config"
	"wannabe/curl/actions"
)

func GenerateCurl(method string, path string, queries map[string]string, headers map[string][]string, body []byte, config config.Config) (string, error) {
	httpMethod := actions.ProcessHttpMethod(method)

	processedHost, err := actions.ProcessHost(config.Server, config.RequestMatching.Host)
	if err != nil {
		return "", err
	}

	processedPath, err := actions.ProcessPath(path, config.RequestMatching.Path)
	if err != nil {
		return "", err
	}

	processedQuery, err := actions.ProcessQuery(queries, config.RequestMatching.Query)
	if err != nil {
		return "", err
	}

	processedHeaders := actions.ProcessHeaders(headers, config.RequestMatching.Headers)

	processedBody, err := actions.ProcessBody(body, config.RequestMatching.Body)
	if err != nil {
		return "", err
	}

	url := actions.PrepareUrl(processedHost, processedPath, processedQuery)

	request, err := actions.PrepareRequest(httpMethod, url, processedBody)
	if err != nil {
		return "", err
	}

	actions.SetHeaders(request, processedHeaders)

	curl, err := actions.PrepareCurl(request)
	if err != nil {
		return "", err
	}

	return curl, nil
}
