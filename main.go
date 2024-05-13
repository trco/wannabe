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

// FIXME
// when request is proxied from IC container connection is not closed after response is received
// github issue opened: https://github.com/AdguardTeam/gomitmproxy/issues/27
func main() {
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

	proxy.Close()
}
