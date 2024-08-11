package handlers

import (
	"fmt"
	"net/http"
	curl "wannabe/curl/services"
	"wannabe/handlers/utils"
	"wannabe/hash/actions"
	"wannabe/providers"
	requestActions "wannabe/request/actions"
	requestServices "wannabe/request/services"
	"wannabe/response/services"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(config types.Config, storageProvider providers.StorageProvider) types.WannabeOnRequestHandler {
	return func(session *gomitmproxy.Session) (*http.Request, *http.Response) {
		return processSessionOnRequest(config, storageProvider, session)
	}
}

func processSessionOnRequest(config types.Config, storageProvider providers.StorageProvider, session *gomitmproxy.Session) (*http.Request, *http.Response) {
	request := session.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil, nil
	}

	request = requestServices.ProcessRequest(request)

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	curl, err := curl.GenerateCurl(request, wannabe)
	if err != nil {
		return utils.InternalErrorOnRequest(session, request, err)
	}
	session.SetProp("curl", curl)

	hash, err := actions.GenerateHash(curl)
	if err != nil {
		return utils.InternalErrorOnRequest(session, request, err)
	}
	session.SetProp("hash", hash)

	isNotProxyMode := config.Mode != types.ProxyMode
	if isNotProxyMode {
		records, err := storageProvider.ReadRecords(host, []string{hash})
		if err != nil {
			return utils.InternalErrorOnRequest(session, request, err)
		}

		isSingleRecord := len(records) == 1
		if isSingleRecord {
			return processRecords(session, request, records[0])
		}

		isServerMode := config.Mode == types.ServerMode
		if isServerMode {
			return utils.InternalErrorOnRequest(session, request, fmt.Errorf("no record found for the request"))
		}
	}

	hasBody := request.Body != nil
	if hasBody {
		requestBody, err := requestActions.CopyBody(request)
		if err != nil {
			return utils.InternalErrorOnRequest(session, request, err)
		}
		session.SetProp("requestBody", requestBody)
	}

	return nil, nil
}

func processRecords(session *gomitmproxy.Session, request *http.Request, record []byte) (*http.Request, *http.Response) {
	responseSetFromRecord, err := services.SetResponse(record, request)
	if err != nil {
		return utils.InternalErrorOnRequest(session, request, err)
	}

	session.SetProp("responseSetFromRecord", true)

	fmt.Println("Response successfully read from configured StorageProvider.")

	return nil, responseSetFromRecord
}
