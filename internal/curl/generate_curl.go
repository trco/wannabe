package curl

import (
	"net/http"

	"github.com/trco/wannabe/internal/config"
)

func generateCurl(request *http.Request, wannabe config.Wannabe) (string, error) {
	payload, err := generateCurlPayload(request)
	if err != nil {
		return "", err
	}

	httpMethod := processHttpMethod(payload.httpMethod)

	processedHost, err := processHost(payload.host, wannabe.RequestMatching.Host)
	if err != nil {
		return "", err
	}

	processedPath, err := processPath(payload.path, wannabe.RequestMatching.Path)
	if err != nil {
		return "", err
	}

	processedQuery, err := processQuery(payload.query, wannabe.RequestMatching.Query)
	if err != nil {
		return "", err
	}

	processedHeaders := processHeaders(payload.requestHeaders, wannabe.RequestMatching.Headers)

	processedBody, err := processBody(payload.requestBody, wannabe.RequestMatching.Body)
	if err != nil {
		return "", err
	}

	url := prepareUrl(payload.scheme, processedHost, processedPath, processedQuery)

	processedRequest, err := prepareRequest(httpMethod, url, processedBody)
	if err != nil {
		return "", err
	}

	setHeaders(processedRequest, processedHeaders)

	curl, err := prepareCurl(processedRequest)
	if err != nil {
		return "", err
	}

	return curl, nil
}
