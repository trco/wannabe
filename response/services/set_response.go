package services

import (
	"net/http"
	"wannabe/response/actions"
)

func SetResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	return actions.SetResponse(encodedRecord, request)
}
