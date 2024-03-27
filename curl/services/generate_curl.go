package services

import (
	"wannabe/config"
	"wannabe/curl/actions"
	"wannabe/curl/entities"
)

func GenerateCurl(config config.Config, payload entities.GenerateCurlPayload) (string, error) {
	httpMethod := actions.ProcessHttpMethod(payload.HttpMethod)

	processedHost, err := actions.ProcessHost(config.Server, config.RequestMatching.Host)
	if err != nil {
		return "", err
	}

	processedPath, err := actions.ProcessPath(payload.Path, config.RequestMatching.Path)
	if err != nil {
		return "", err
	}

	processedQuery, err := actions.ProcessQuery(payload.Query, config.RequestMatching.Query)
	if err != nil {
		return "", err
	}

	processedHeaders := actions.ProcessHeaders(payload.RequestHeaders, config.RequestMatching.Headers)

	processedBody, err := actions.ProcessBody(payload.RequestBody, config.RequestMatching.Body)
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
