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
		request := session.Request()

		if request.Method == "CONNECT" {
			return nil
		}

		if shouldSkipResponseProcessing(session) {
			return nil
		}

		hash, curl, err := getHashAndCurlFromSession(session)
		if err != nil {
			return internalErrorOnResponse(session, request, err)
		}

		recordPayload, err := record.GenerateRecordPayload(session, hash, curl)
		if err != nil {
			return internalErrorOnResponse(session, request, err)
		}

		host := request.URL.Host
		wannabe := config.Wannabes[host]

		record, err := record.GenerateRecord(wannabe.Records, recordPayload)
		if err != nil {
			return internalErrorOnResponse(session, request, err)
		}

		err = storageProvider.InsertRecords([][]byte{record}, []string{hash}, host)
		if err != nil {
			return internalErrorOnResponse(session, request, err)
		}

		return nil
	}
}
