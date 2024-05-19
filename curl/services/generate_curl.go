package services

import (
	"net/http"
	"wannabe/curl/actions"
	"wannabe/types"
)

func GenerateCurl(request *http.Request, wannabe types.Wannabe) (string, error) {
	payload, err := actions.GenerateCurlPayload(request)
	if err != nil {
		return "", err
	}

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

	processedRequest, err := actions.PrepareRequest(httpMethod, url, processedBody)
	if err != nil {
		return "", err
	}

	actions.SetHeaders(processedRequest, processedHeaders)

	curl, err := actions.PrepareCurl(processedRequest)
	if err != nil {
		return "", err
	}

	return curl, nil
}
