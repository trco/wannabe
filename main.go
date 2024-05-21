package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	cfg "wannabe/config"
	"wannabe/handlers"
	"wannabe/providers"
	"wannabe/types"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/mitm"
)

func main() {
	config, err := cfg.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	mitmConfig, err := cfg.LoadMitmConfig("demo.crt", "demo.key")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := providers.StorageProviderFactory(config)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	go startWannabeApiServer(config, storageProvider)
	startWannabeProxyServer(config, mitmConfig, storageProvider)
}

func startWannabeProxyServer(config types.Config, mitmConfig *mitm.Config, storageProvider providers.StorageProvider) {
	proxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 6789,
		},
		MITMConfig: mitmConfig,
		OnRequest:  handlers.WannabeOnRequest(config, storageProvider),
		OnResponse: handlers.WannabeOnResponse(config, storageProvider),
	})

	err := proxy.Start()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	proxy.Close()
}

func startWannabeApiServer(config types.Config, storageProvider providers.StorageProvider) {
	http.HandleFunc("/wannabe/api/records/{hash}", handlers.Records(config, storageProvider))
	http.HandleFunc("/wannabe/api/records", handlers.Records(config, storageProvider))

	fmt.Println("Listening on localhost:1234")
	err := http.ListenAndServe("localhost:1234", nil)
	if err != nil {
		log.Fatal(err)
	}
}
