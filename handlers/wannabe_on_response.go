package handlers

import (
	"net/http"
	"wannabe/handlers/utils"
	"wannabe/providers"
	record "wannabe/record/services"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(config types.Config, storageProvider providers.StorageProvider) types.WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		return processSessionOnResponse(config, storageProvider, session)
	}
}

func processSessionOnResponse(config types.Config, storageProvider providers.StorageProvider, session *gomitmproxy.Session) *http.Response {
	request := session.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil
	}

	if utils.ShouldSkipResponseProcessing(session) {
		return nil
	}

	hash, curl, err := utils.GetHashAndCurlFromSession(session)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	recordPayload, err := record.GenerateRecordPayload(session, hash, curl)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	record, err := record.GenerateRecord(wannabe.Records, recordPayload)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{record}, false)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	return nil
}
