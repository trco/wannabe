package handlers

import (
	"fmt"
	"net/http"
	cfg "wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	response "wannabe/response/services"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(config cfg.Config, storageProvider providers.StorageProvider) WannabeOnRequestHandler {
	return func(session *gomitmproxy.Session) (*http.Request, *http.Response) {
		return processSessionOnRequest(config, storageProvider, session)
	}
}

func processSessionOnRequest(config cfg.Config, storageProvider providers.StorageProvider, session *gomitmproxy.Session) (*http.Request, *http.Response) {
	request := session.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil, nil
	}

	host := request.URL.Host
	wannabe := config.Wannabes[host]

	curl, err := curl.GenerateCurl(request, wannabe)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}
	session.SetProp("curl", curl)

	hash, err := hash.GenerateHash(curl)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}
	session.SetProp("hash", hash)

	isNotProxyMode := config.Mode != cfg.ProxyMode
	if isNotProxyMode {
		records, err := storageProvider.ReadRecords(host, []string{hash})
		if err != nil {
			return internalErrorOnRequest(session, request, err)
		}

		isSingleRecord := len(records) == 1
		if isSingleRecord {
			return processRecords(session, request, records[0])
		}

		isServerMode := config.Mode == cfg.ServerMode
		if isServerMode {
			return internalErrorOnRequest(session, request, fmt.Errorf("no record found for the request"))
		}
	}

	return nil, nil
}

func processRecords(session *gomitmproxy.Session, request *http.Request, record []byte) (*http.Request, *http.Response) {
	responseSetFromRecord, err := response.SetResponse(record, request)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}

	session.SetProp("responseSetFromRecord", true)

	// TODO remove log
	fmt.Println("GetResponse >>> READ and return")

	return nil, responseSetFromRecord
}
