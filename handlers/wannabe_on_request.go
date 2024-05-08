package handlers

import (
	"log"
	"net/http"
	"wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"

	"github.com/AdguardTeam/gomitmproxy"
)

func WannabeOnRequest(configuration config.Config, storageProvider providers.StorageProvider) WannabeHandler {
	return func(session *gomitmproxy.Session) (request *http.Request, response *http.Response) {
		originalRequest := session.Request()

		if originalRequest.Method == "CONNECT" {
			return nil, nil
		}

		host := originalRequest.URL.Host
		wannabe := configuration.Wannabes[host]

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
		if configuration.Mode != "proxy" {
			// records, err := storageProvider.ReadRecords([]string{hash}, host)
			_, err := storageProvider.ReadRecords([]string{hash}, host)
			if err != nil && configuration.FailOnReadError {
				return internalError(session, originalRequest, err)
			}

			// if records != nil {
			// 	// response from record is set directly to ctx
			// 	res, err = response.SetResponse(ctx, records[0])
			// 	if err != nil && config.FailOnReadError {
			// 		return internalError(session, req, err)
			// 	}

			// 	// TODO remove log
			// 	fmt.Println("GetResponse >>> READ and return")

			// 	// FIXME probably not OK
			// 	// FIXME
			// 	// 1. add info to session.props that response was prepared from record
			// 	// 2. read prop in OnResponse handler and simply return response prepared from record
			//  // 3. see "session.SetProp("blocked", true)" in internal error
			// 	return nil, res
			// }

			// if config.Mode == "server" {
			// 	return internalError(session, req, fmt.Errorf("no record found for the request"))
			// }
			return nil, nil
		}

		return nil, nil
	}
}
