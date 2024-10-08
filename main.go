package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/mitm"
	"github.com/trco/wannabe/internal/config"
	"github.com/trco/wannabe/internal/handler"
	"github.com/trco/wannabe/internal/storage"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	storageProvider, err := storage.ProviderFactory(cfg.StorageProvider)
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	mitmConfig, err := config.LoadMitmConfig()
	if err != nil {
		log.Fatalf("fatal error starting app: %v", err)
	}

	go startWannabeApiServer(cfg, storageProvider)
	startWannabeProxyServer(cfg, mitmConfig, storageProvider)
}

func startWannabeApiServer(cfg config.Config, storageProvider storage.Provider) {
	http.HandleFunc("/wannabe/api/records/{hash}", handler.Records(cfg, storageProvider))
	http.HandleFunc("/wannabe/api/records", handler.Records(cfg, storageProvider))
	http.HandleFunc("/wannabe/api/regenerate", handler.Regenerate(cfg, storageProvider))

	fmt.Println("API server start listening to [::]:6790")
	err := http.ListenAndServe(":6790", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func startWannabeProxyServer(cfg config.Config, mitmConfig *mitm.Config, storageProvider storage.Provider) {
	wannabeProxy := initWannabeProxy(cfg, mitmConfig, storageProvider)

	err := wannabeProxy.Start()
	if err != nil {
		log.Fatal(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChannel

	wannabeProxy.Close()
}

func initWannabeProxy(cfg config.Config, mitmConfig *mitm.Config, storageProvider storage.Provider) *gomitmproxy.Proxy {
	wannabeProxy := gomitmproxy.NewProxy(gomitmproxy.Config{
		ListenAddr: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 6789,
		},
		MITMConfig: mitmConfig,
		OnRequest:  handler.Request(cfg, storageProvider),
		OnResponse: handler.Response(cfg, storageProvider),
	})

	return wannabeProxy
}
