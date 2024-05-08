package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

func internalError(session *gomitmproxy.Session, originalRequest *http.Request, err error) (request *http.Request, response *http.Response) {
	body := prepareResponseBody(err)
	res := proxyutil.NewResponse(http.StatusInternalServerError, body, originalRequest)
	res.Header.Set("Content-Type", "text/html")

	// REVIEW and use session.Props
	// Use session props to pass the information about request being blocked
	session.SetProp("blocked", true)
	return nil, res
}

func prepareResponseBody(err error) *bytes.Reader {
	body, err := json.Marshal(InternalError{
		Error: err.Error(),
	})
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	bodyReader := bytes.NewReader(body)

	return bodyReader
}

func checkDuplicates(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

func processRecordValidation(recordProcessingDetails *[]RecordProcessingDetails, hash string, message string, valueToIncrement *int) {
	*recordProcessingDetails = append(*recordProcessingDetails, RecordProcessingDetails{
		Hash:    hash,
		Message: message,
	})

	*valueToIncrement++
}
