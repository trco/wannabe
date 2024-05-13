package handlers

import (
	"net/http"

	"github.com/AdguardTeam/gomitmproxy"
)

type WannabeOnRequestHandler func(session *gomitmproxy.Session) (*http.Request, *http.Response)
type WannabeOnResponseHandler func(session *gomitmproxy.Session) *http.Response

type InternalError struct {
	Error string `json:"error"`
}

type RecordProcessingDetails struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

type PostRecordsResponse struct {
	InsertedRecordsCount    int                       `json:"insertedRecordsCount"`
	NotInsertedRecordsCount int                       `json:"notInsertedRecordsCount"`
	RecordProcessingDetails []RecordProcessingDetails `json:"recordProcessingDetails"`
}

type DeleteRecordsResponse struct {
	Message string   `json:"message"`
	Hashes  []string `json:"hashes"`
}

type RegenerateResponse struct {
	Message           string   `json:"message"`
	RegeneratedHashes []string `json:"regeneratedHashes"`
	FailedHashes      []string `json:"failedHashes"`
}
