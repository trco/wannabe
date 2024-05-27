package types

import (
	"net/http"

	"github.com/AdguardTeam/gomitmproxy"
)

type WannabeOnRequestHandler func(*gomitmproxy.Session) (*http.Request, *http.Response)
type WannabeOnResponseHandler func(*gomitmproxy.Session) *http.Response

type WannabeSession struct {
	Req   *http.Request
	Res   *http.Response
	Props map[string]interface{}
}

func (s *WannabeSession) GetRequest() *http.Request {
	return s.Req
}

func (s *WannabeSession) GetResponse() *http.Response {
	return s.Res
}

func (s *WannabeSession) GetProp(key string) (interface{}, bool) {
	v, ok := s.Props[key]
	return v, ok
}

func (s *WannabeSession) SetProp(key string, val interface{}) {
	s.Props[key] = val
}

type InternalError struct {
	Error string `json:"error"`
}

type InternalErrorApi struct {
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
