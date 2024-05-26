package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

// wannabe

func internalErrorOnRequest(wannabeSession types.WannabeSession, request *http.Request, err error) (*http.Request, *http.Response) {
	wannabeSession.SetProp("blocked", true)

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

func shouldSkipResponseProcessing(wannabeSession types.WannabeSession) bool {
	if _, blocked := wannabeSession.GetProp("blocked"); blocked {
		return true
	}
	if _, responseSetFromRecord := wannabeSession.GetProp("responseSetFromRecord"); responseSetFromRecord {
		return true
	}
	return false
}

func getHashAndCurlFromSession(wannabeSession types.WannabeSession) (string, string, error) {
	hashProp, ok := wannabeSession.GetProp("hash")
	if !ok {
		return "", "", fmt.Errorf("no hash in session")
	}
	hash, ok := hashProp.(string)
	if !ok {
		return "", "", fmt.Errorf("hash is not a string")
	}

	curlProp, ok := wannabeSession.GetProp("curl")
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
