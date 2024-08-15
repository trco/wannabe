package services

import (
	"net/http"

	"github.com/trco/wannabe/request/actions"
)

func ProcessRequest(request *http.Request) *http.Request {
	processedRequest := actions.RemoveBody(request)
	processedRequest = actions.ProcessScheme(processedRequest)

	return processedRequest
}
