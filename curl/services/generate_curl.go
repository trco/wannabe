package services

import (
	"wannabe/config"
	"wannabe/curl/actions"
	"wannabe/curl/entities"
)

func GenerateCurl(payload entities.GenerateCurlPayload, wannabe config.Wannabe) (string, error) {
	httpMethod := actions.ProcessHttpMethod(payload.HttpMethod)

	processedHost, err := actions.ProcessHost(payload.Host, wannabe.RequestMatching.Host)
	if err != nil {
		return "", err
	}

	processedPath, err := actions.ProcessPath(payload.Path, wannabe.RequestMatching.Path)
	if err != nil {
		return "", err
	}

	processedQuery, err := actions.ProcessQuery(payload.Query, wannabe.RequestMatching.Query)
	if err != nil {
		return "", err
	}

	processedHeaders := actions.ProcessHeaders(payload.RequestHeaders, wannabe.RequestMatching.Headers)

	processedBody, err := actions.ProcessBody(payload.RequestBody, wannabe.RequestMatching.Body)
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
