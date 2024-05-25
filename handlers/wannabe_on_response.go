package handlers

import (
	"net/http"
	"wannabe/providers"
	record "wannabe/record/services"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(config types.Config, storageProvider providers.StorageProvider) types.WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		wannabeSession := types.WannabeSession{
			Req: session.Request(),
			Res: session.Response(),
		}
		return processSessionOnResponse(config, storageProvider, wannabeSession)
	}
}

func processSessionOnResponse(config types.Config, storageProvider providers.StorageProvider, wannabeSession types.WannabeSession) *http.Response {
	request := wannabeSession.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil
	}

	if shouldSkipResponseProcessing(wannabeSession) {
		return nil
	}

	hash, curl, err := getHashAndCurlFromSession(wannabeSession)
	if err != nil {
		return wannabeOnResponseInternalError(request, err)
	}

	recordPayload, err := record.GenerateRecordPayload(wannabeSession, hash, curl)
	if err != nil {
		return wannabeOnResponseInternalError(request, err)
	}

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	record, err := record.GenerateRecord(wannabe.Records, recordPayload)
	if err != nil {
		return wannabeOnResponseInternalError(request, err)
	}

	err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{record})
	if err != nil {
		return wannabeOnResponseInternalError(request, err)
	}

	return nil
}
