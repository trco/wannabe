package handlers

import (
	"fmt"
	"log"
	"net/http"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"
	response "wannabe/response/services"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(config config.Config, storageProvider providers.StorageProvider) WannabeOnRequestHandler {
	return func(session *gomitmproxy.Session) (wannabeRequest *http.Request, wannabeResponse *http.Response) {
		originalRequest := session.Request()

		if originalRequest.Method == "CONNECT" {
			return nil, nil
		}

		host := originalRequest.URL.Host
		wannabe := config.Wannabes[host]

		log.Printf("originalRequest:", originalRequest)
		log.Printf("onRequest: %s %s %s", originalRequest.Method, originalRequest.URL.String(), originalRequest.URL.Host)
		log.Printf("host: %s", host)

		curl, err := curl.GenerateCurl(originalRequest, wannabe)
		log.Printf("curl: %s", curl)
		if err != nil {
			return internalError(session, originalRequest, err)
		}

		hash, err := hash.GenerateHash(curl)
		log.Printf("hash: %s", hash)
		if err != nil {
			return internalError(session, originalRequest, err)
		}

		// server, mixed
		if config.Mode != "proxy" {
			records, err := storageProvider.ReadRecords([]string{hash}, host)
			if err != nil && config.FailOnReadError {
				return internalError(session, originalRequest, err)
			}

			if records != nil {
				wannabeResponse, err := response.SetResponse(records[0], originalRequest)
				if err != nil && config.FailOnReadError {
					return internalError(session, originalRequest, err)
				}

				// TODO remove log
				fmt.Println("GetResponse >>> READ and return")

				// FIXME probably not OK
				// FIXME
				// 1. add info to session.props that response was prepared from record
				// 2. read prop in OnResponse handler and simply return response prepared from record
				// 3. see "session.SetProp("blocked", true)" in internal error
				session.SetProp("responseReadFromStorageProvider", true)

				return nil, wannabeResponse
			}

			if config.Mode == "server" {
				session.SetProp("blocked", true)
				return internalError(session, originalRequest, fmt.Errorf("no record found for the request"))
			}

			return nil, nil
		}

		return nil, nil
	}
}
