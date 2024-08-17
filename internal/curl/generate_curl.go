package curl

import (
	"net/http"

	"github.com/trco/wannabe/internal/config"
)

func GenerateCurl(request *http.Request, wannabe config.Wannabe) (string, error) {
	payload, err := GenerateCurlPayload(request)
	if err != nil {
		return "", err
	}

	httpMethod := ProcessHttpMethod(payload.HttpMethod)

	processedHost, err := ProcessHost(payload.Host, wannabe.RequestMatching.Host)
	if err != nil {
		return "", err
	}

	processedPath, err := ProcessPath(payload.Path, wannabe.RequestMatching.Path)
	if err != nil {
		return "", err
	}

	processedQuery, err := ProcessQuery(payload.Query, wannabe.RequestMatching.Query)
	if err != nil {
		return "", err
	}

	processedHeaders := ProcessHeaders(payload.RequestHeaders, wannabe.RequestMatching.Headers)

	processedBody, err := ProcessBody(payload.RequestBody, wannabe.RequestMatching.Body)
	if err != nil {
		return "", err
	}

	url := PrepareUrl(payload.Scheme, processedHost, processedPath, processedQuery)

	processedRequest, err := PrepareRequest(httpMethod, url, processedBody)
	if err != nil {
		return "", err
	}

	SetHeaders(processedRequest, processedHeaders)

	curl, err := PrepareCurl(processedRequest)
	if err != nil {
		return "", err
	}

	return curl, nil
}
