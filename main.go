package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/trco/wannabe/internal/handlers"
	"github.com/trco/wannabe/internal/providers"
	"github.com/trco/wannabe/internal/types"

	cfg "github.com/trco/wannabe/internal/config"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/mitm"
)

func main() {
	config, err := cfg.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := providers.StorageProviderFactory(config)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	mitmConfig, err := cfg.LoadMitmConfig("wannabe.crt", "wannabe.key")
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	go startWannabeApiServer(config, storageProvider)
	startWannabeProxyServer(config, mitmConfig, storageProvider)
}

func startWannabeApiServer(config types.Config, storageProvider providers.StorageProvider) {
	http.HandleFunc("/wannabe/api/records/{hash}", handlers.Records(config, storageProvider))
	http.HandleFunc("/wannabe/api/records", handlers.Records(config, storageProvider))
	http.HandleFunc("/wannabe/api/regenerate", handlers.Regenerate(config, storageProvider))

	fmt.Println("API server start listening to [::]:6790")
	err := http.ListenAndServe(":6790", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startWannabeProxyServer(config types.Config, mitmConfig *mitm.Config, storageProvider providers.StorageProvider) {
	wannabeProxy := initWannabeProxy(config, mitmConfig, storageProvider)

	err := wannabeProxy.Start()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	wannabeProxy.Close()
}

func initWannabeProxy(config types.Config, mitmConfig *mitm.Config, storageProvider providers.StorageProvider) *gomitmproxy.Proxy {
	wannabeProxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 6789,
		},
		MITMConfig: mitmConfig,
		OnRequest:  handlers.WannabeOnRequest(config, storageProvider),
		OnResponse: handlers.WannabeOnResponse(config, storageProvider),
	})

	return wannabeProxy
}
