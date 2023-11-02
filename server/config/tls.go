package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

type tlsOptionFunction func(*tls.Config) error

type TLSBuilder struct {
	CertProvider CertificateProvider
}

func UseTLSVersion(version uint16) tlsOptionFunction {
	return func(config *tls.Config) error {
		if version < tls.VersionTLS10 || version > tls.VersionTLS13 {
			return errors.New("Invalid or unsafe TLS/SSL version provided")
		}
		config.MinVersion = version
		config.MaxVersion = version
		return nil
	}
}

func useMutualTLS(config *tls.Config) error {
	config.ClientCAs = config.RootCAs
	config.ClientAuth = tls.RequireAndVerifyClientCert
	return nil
}

func withRootCA(rootCACertPool *x509.CertPool) tlsOptionFunction {
	return func(config *tls.Config) error {
		config.RootCAs = rootCACertPool
		return nil
	}
}

func withCertificiates(certs []tls.Certificate) tlsOptionFunction {
	return func(config *tls.Config) error {
		config.Certificates = certs
		return nil
	}
}

func defaultTLSConfig() *tls.Config {
	return &tls.Config{MinVersion: tls.VersionTLS12}
}

func (b *TLSBuilder) BuildTLS(tlsOptions ...tlsOptionFunction) (*tls.Config, error) {
	tlsConfig := defaultTLSConfig()
	certs, err := b.CertProvider.LoadServerCertificates()
	if err != nil {
		return nil, err
	}
	withCertificiates(certs)(tlsConfig)

	rootCertPool, err := b.CertProvider.LoadRootCertificatePool()
	if err != nil {
		return nil, err
	}
	withRootCA(rootCertPool)(tlsConfig)

	for _, function := range tlsOptions {
		if err := function(tlsConfig); err != nil {
			return nil, err
		}
	}
	return tlsConfig, nil
}
