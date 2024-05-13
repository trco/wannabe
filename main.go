package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	configPackage "wannabe/config"
	"wannabe/handlers"
	"wannabe/providers"

	"github.com/AdguardTeam/gomitmproxy"
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
	config, err := configPackage.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := providers.StorageProviderFactory(config)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	mitmConfig, err := configPackage.LoadMitmConfig("demo.crt", "demo.key")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	proxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 6789,
		},
		MITMConfig: mitmConfig,
		OnRequest:  handlers.WannabeOnRequest(config, storageProvider),
		OnResponse: handlers.WannabeOnResponse(config, storageProvider),
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
