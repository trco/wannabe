package services

import (
	"net/http"
	"wannabe/response/actions"
)

func SetResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	response, err := actions.SetResponse(encodedRecord, request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
