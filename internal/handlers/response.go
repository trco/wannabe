package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/trco/wannabe/internal/config"
	"github.com/trco/wannabe/internal/record"
	"github.com/trco/wannabe/internal/storage"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

type ResponseHandler func(*gomitmproxy.Session) *http.Response

func Response(cfg config.Config, storageProvider storage.Provider) ResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		return processResponse(cfg, storageProvider, session)
	}
}

func processResponse(cfg config.Config, storageProvider storage.Provider, session *gomitmproxy.Session) *http.Response {
	request := session.Request()

	if requestBody, ok := session.GetProp("requestBody"); ok {
		request.Body = requestBody.(io.ReadCloser)
	}

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil
	}

	if shouldSkipResponseProcessing(session) {
		return nil
	}

	hash, curl, err := getHashAndCurlFromSession(session)
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	recordPayload, err := record.GenerateRecordPayload(session, hash, curl)
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	host := request.URL.Host
	wannabe := cfg.Wannabes[host]

	record, err := record.GenerateRecord(wannabe.Records, recordPayload)
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{record}, false)
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	return nil
}

func shouldSkipResponseProcessing(session Session) bool {
	if _, blocked := session.GetProp("blocked"); blocked {
		return true
	}
	if _, response := session.GetProp("responseSet"); response {
		return true
	}
	return false
}

func getHashAndCurlFromSession(session Session) (string, string, error) {
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

func internalErrorOnResponse(request *http.Request, err error) *http.Response {
	body := PrepareResponseBody(err)
	response := proxyutil.NewResponse(http.StatusInternalServerError, body, request)
	response.Header.Set("Content-Type", "application/json")

	return response
}
