package handlers

import (
	"fmt"
	"net/http"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	response "wannabe/response/services"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(config config.Config, storageProvider providers.StorageProvider) WannabeOnRequestHandler {
	return func(session *gomitmproxy.Session) (*http.Request, *http.Response) {
		request := session.Request()

		if request.Method == "CONNECT" {
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

		// server, mixed
		if config.Mode != "proxy" {
			records, err := storageProvider.ReadRecords([]string{hash}, host)
			if err != nil && config.FailOnReadError {
				return internalErrorOnRequest(session, request, err)
			}

			if records != nil {
				return processRecords(session, request, records, config.FailOnReadError)
			}

			if config.Mode == "server" {
				return internalErrorOnRequest(session, request, fmt.Errorf("no record found for the request"))
			}
		}

		return nil, nil
	}
}

func processRecords(session *gomitmproxy.Session, request *http.Request, records [][]byte, failOnReadError bool) (*http.Request, *http.Response) {
	responseFromRecord, err := response.SetResponse(records[0], request)
	if err != nil && failOnReadError {
		return internalErrorOnRequest(session, request, err)
	}

	session.SetProp("responseFromRecord", true)

	// TODO remove log
	fmt.Println("GetResponse >>> READ and return")

	return nil, responseFromRecord
}
