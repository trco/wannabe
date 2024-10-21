package config

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/AdguardTeam/gomitmproxy/mitm"
)

func LoadMitmConfig() (*mitm.Config, error) {
	certPath, certKeyPath, err := getCertPaths()
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.LoadX509KeyPair(certPath, certKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey := tlsCert.PrivateKey.(*rsa.PrivateKey)

	x509c, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, err
	}

	mitmConfig, err := mitm.NewConfig(x509c, privateKey, nil)
	if err != nil {
		return nil, err
	}

	return mitmConfig, nil
}

func getCertPaths() (string, string, error) {
	certPath := os.Getenv(CertPath)
	if certPath == "" {
		certPath = "certs/wannabe.crt"

		if _, err := os.Stat(certPath); err != nil {
			return "", "", fmt.Errorf("failed loading certificate file: %v", err)
		}
	}

	certKeyPath := os.Getenv(CertKeyPath)
	if certKeyPath == "" {
		certKeyPath = "certs/wannabe.key"

		if _, err := os.Stat(certKeyPath); err != nil {
			return "", "", fmt.Errorf("failed loading key file: %v", err)
		}
	}

	return certPath, certKeyPath, nil
}
