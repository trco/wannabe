package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	config "wannabe/config"
	curl "wannabe/curl/services"
	hash "wannabe/hash/services"
	"wannabe/providers"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

// 1. create self-signed certificate
// openssl genrsa -out demo.key 2048
// openssl req -new -x509 -key demo.key -out demo.crt -days 3650

// 2. add it to client container
// docker cp ./demo.crt integrations-core.local:/usr/local/share/ca-certificates/

// 3. enter client container and add certificate to ca-certificates & check that it was added
// update-ca-certificates
// cat /etc/ssl/certs/ca-certificates.crt

// FIXME
// !!! when request is proxied from IC container it doesn't exit after response is returned
// curl -x http://host.docker.internal:6789 https://api.github.com
// !!! this works as expected with basic example, see go/wannabe-proxy project
// curl -x http://host.docker.internal:6667 https://api.github.com

func main() {
	// TODO read config path from env variable
	configuration, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := providers.StorageProviderFactory(configuration)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	mitmConfig, err := config.LoadMitmConfig("demo.crt", "demo.key")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	proxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 6789,
		},
		MITMConfig: mitmConfig,
		OnRequest: func(session *gomitmproxy.Session) (request *http.Request, response *http.Response) {
			originalRequest := session.Request()

			if originalRequest.Method != "CONNECT" {
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
				}
			}

			return nil, nil
		},
		// OnResponse: func(session *gomitmproxy.Session) *http.Response {
		// 	log.Printf("onResponse: %s", session.Request().URL.String())

		// 	if _, ok := session.GetProp("blocked"); ok {
		// 		log.Printf("onResponse: was blocked")
		// 	}

		// 	return nil
		// },
	})
	err = proxy.Start()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	// Clean up
	proxy.Close()
}

func internalError(session *gomitmproxy.Session, req *http.Request, err error) (request *http.Request, response *http.Response) {
	body := prepareResponseBody(err)
	// REVIEW
	// return proxyutil.NewErrorResponse(session.Request(), err)
	res := proxyutil.NewResponse(http.StatusInternalServerError, body, req)
	res.Header.Set("Content-Type", "text/html")

	// Use session props to pass the information about request being blocked
	session.SetProp("blocked", true)
	return nil, res
}

func prepareResponseBody(err error) *bytes.Reader {
	body, err := json.Marshal(InternalError{
		Error: err.Error(),
	})
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	bodyReader := bytes.NewReader(body)

	return bodyReader
}

type InternalError struct {
	Error string `json:"error"`
}
