package handlers

import (
	"fmt"
	"net/http"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	response "wannabe/response/services"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(config types.Config, storageProvider providers.StorageProvider) types.WannabeOnRequestHandler {
	return func(session *gomitmproxy.Session) (*http.Request, *http.Response) {
		wannabeSession := types.WannabeSession{
			Req: session.Request(),
			Res: session.Response(),
		}
		return processSessionOnRequest(config, storageProvider, wannabeSession)
	}
}

func processSessionOnRequest(config types.Config, storageProvider providers.StorageProvider, wannabeSession types.WannabeSession) (*http.Request, *http.Response) {
	request := wannabeSession.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil, nil
	}

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	curl, err := curl.GenerateCurl(request, wannabe)
	if err != nil {
		return internalErrorOnRequest(wannabeSession, request, err)
	}
	wannabeSession.SetProp("curl", curl)

	hash, err := hash.GenerateHash(curl)
	if err != nil {
		return internalErrorOnRequest(wannabeSession, request, err)
	}
	wannabeSession.SetProp("hash", hash)

	isNotProxyMode := config.Mode != types.ProxyMode
	if isNotProxyMode {
		records, err := storageProvider.ReadRecords(host, []string{hash})
		if err != nil {
			return internalErrorOnRequest(wannabeSession, request, err)
		}

		isSingleRecord := len(records) == 1
		if isSingleRecord {
			return processRecords(wannabeSession, request, records[0])
		}

		isServerMode := config.Mode == types.ServerMode
		if isServerMode {
			return internalErrorOnRequest(wannabeSession, request, fmt.Errorf("no record found for the request"))
		}
	}

	return nil, nil
}

func processRecords(wannabeSession types.WannabeSession, request *http.Request, record []byte) (*http.Request, *http.Response) {
	responseSetFromRecord, err := response.SetResponse(record, request)
	if err != nil {
		return internalErrorOnRequest(wannabeSession, request, err)
	}

	wannabeSession.SetProp("responseSetFromRecord", true)

	// TODO remove log
	fmt.Println("GetResponse >>> READ and return")

	return nil, responseSetFromRecord
}
