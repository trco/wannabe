package handlers

import (
	"net/http"
	"wannabe/config"
	"wannabe/providers"
	record "wannabe/record/services"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(config config.Config, storageProvider providers.StorageProvider) WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		return processSessionOnResponse(config, storageProvider, session)
	}
}

func processSessionOnResponse(config config.Config, storageProvider providers.StorageProvider, session *gomitmproxy.Session) *http.Response {
	request := session.Request()

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
	wannabe := config.Wannabes[host]

	record, err := record.GenerateRecord(wannabe.Records, recordPayload)
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{record})
	if err != nil {
		return internalErrorOnResponse(request, err)
	}

	return nil
}
