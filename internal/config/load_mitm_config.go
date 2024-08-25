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
	// var certPath, certKeyPath string
	certPath, certKeyPath, err := getCertPaths()

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
	var certPath, certKeyPath string

	if os.Getenv(RunningInContainer) == "" {
		_, err := os.Stat("wannabe.crt")
		if err != nil {
			return "", "", fmt.Errorf("failed loading wannabe.crt file: %v", err)
		}
		certPath = "wannabe.crt"

		_, err = os.Stat("wannabe.key")
		if err != nil {
			return "", "", fmt.Errorf("failed loading wannabe.key file: %v", err)
		}
		certKeyPath = "wannabe.key"

		return certPath, certKeyPath, nil
	}

	certPath = os.Getenv(CertPath)
	if certPath == "" {
		return "", "", fmt.Errorf("%v env variable not set", CertPath)
	}

	certKeyPath = os.Getenv(CertKeyPath)
	if certPath == "" {
		return "", "", fmt.Errorf("%v env variable not set", CertKeyPath)
	}

	return certPath, certKeyPath, nil

}
