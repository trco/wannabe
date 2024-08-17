package handlers

import (
	"io"
	"net/http"

	"github.com/trco/wannabe/handlers/utils"
	"github.com/trco/wannabe/providers"
	"github.com/trco/wannabe/record/actions"
	"github.com/trco/wannabe/record/services"
	"github.com/trco/wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnResponse(config types.Config, storageProvider providers.StorageProvider) types.WannabeOnResponseHandler {
	return func(session *gomitmproxy.Session) *http.Response {
		return processSessionOnResponse(config, storageProvider, session)
	}
}

func processSessionOnResponse(config types.Config, storageProvider providers.StorageProvider, session *gomitmproxy.Session) *http.Response {
	request := session.Request()

	if requestBody, ok := session.GetProp("requestBody"); ok {
		request.Body = requestBody.(io.ReadCloser)
	}

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

	recordPayload, err := actions.GenerateRecordPayload(session, hash, curl)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	record, err := services.GenerateRecord(wannabe.Records, recordPayload)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	err = storageProvider.InsertRecords(host, []string{hash}, [][]byte{record}, false)
	if err != nil {
		return utils.InternalErrorOnResponse(request, err)
	}

	return nil
}
