package services

import (
	"net/http"
	"wannabe/request/actions"
)

func ProcessRequest(request *http.Request) *http.Request {
	processedRequest := actions.ProcessGet(request)
	processedRequest = actions.ProcessHttps(processedRequest)

	return processedRequest
}
