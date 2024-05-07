package config

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/AdguardTeam/gomitmproxy/mitm"
)

func LoadMitmConfig(certificate string, key string) (*mitm.Config, error) {
	// setup mitm config for TLC interception
	// ref: https://github.com/AdguardTeam/gomitmproxy?tab=readme-ov-file#tls-interception
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

	mitmConfig.SetValidity(time.Hour * 24 * 7) // generate certs valid for 7 days
	mitmConfig.SetOrganization("gomitmproxy")  // cert organization

	return mitmConfig, nil
}
