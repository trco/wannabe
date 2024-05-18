package config

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"

	"github.com/AdguardTeam/gomitmproxy/mitm"
)

func LoadMitmConfig(certificate string, key string) (*mitm.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair(certificate, key)
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
