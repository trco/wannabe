package services

import (
	"net/http"
	"wannabe/record/actions"
	"wannabe/types"
)

func GenerateRequest(recordRequest types.Request) (*http.Request, error) {
	return actions.GenerateRequest(recordRequest)
}
