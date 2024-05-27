package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

// wannabe

func internalErrorOnRequest(session *gomitmproxy.Session, request *http.Request, err error) (*http.Request, *http.Response) {
	session.SetProp("blocked", true)

	body := prepareResponseBody(err)
	response := proxyutil.NewResponse(http.StatusInternalServerError, body, request)
	response.Header.Set("Content-Type", "application/json")

	return nil, response
}

func internalErrorOnResponse(request *http.Request, err error) *http.Response {
	body := prepareResponseBody(err)
	response := proxyutil.NewResponse(http.StatusInternalServerError, body, request)
	response.Header.Set("Content-Type", "application/json")

	return response
}

func prepareResponseBody(err error) *bytes.Reader {
	body, err := json.Marshal(types.InternalError{
		Error: err.Error(),
	})
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	bodyReader := bytes.NewReader(body)

	return bodyReader
}

func shouldSkipResponseProcessing(session *gomitmproxy.Session) bool {
	if _, blocked := session.GetProp("blocked"); blocked {
		return true
	}
	if _, responseSetFromRecord := session.GetProp("responseSetFromRecord"); responseSetFromRecord {
		return true
	}
	return false
}

func getHashAndCurlFromSession(session *gomitmproxy.Session) (string, string, error) {
	hashProp, ok := session.GetProp("hash")
	if !ok {
		return "", "", fmt.Errorf("no hash in session")
	}
	hash, ok := hashProp.(string)
	if !ok {
		return "", "", fmt.Errorf("hash is not a string")
	}

	curlProp, ok := session.GetProp("curl")
	if !ok {
		return "", "", fmt.Errorf("no curl in session")
	}
	curl, ok := curlProp.(string)
	if !ok {
		return "", "", fmt.Errorf("curl is not a string")
	}

	return hash, curl, nil
}

// wannabe api

func internalErrorApi(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(types.InternalErrorApi{Error: err.Error()})
}

func apiResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func checkDuplicates(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

func processRecordValidation(recordProcessingDetails *[]types.RecordProcessingDetails, hash string, message string, valueToIncrement *int) {
	*recordProcessingDetails = append(*recordProcessingDetails, types.RecordProcessingDetails{
		Hash:    hash,
		Message: message,
	})

	*valueToIncrement++
}
