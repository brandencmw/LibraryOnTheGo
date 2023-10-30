package config

import (
	"crypto/tls"
)

type TLSConfigProvider interface {
	GetTLSConfig() (*tls.Config, error)
}

type TLS13ConfigProvider struct {
	certProvider    X509CertificateProvider
	leafCertFile    string
	leafCertKeyFile string
	rootCAFiles     []string
}

func NewTLS13ConfigProvider(leafCertFile, leafCertKeyFile string, rootCAFiles []string) *TLS13ConfigProvider {
	return &TLS13ConfigProvider{
		certProvider:    &TLSCertificateProvider{},
		leafCertFile:    leafCertFile,
		leafCertKeyFile: leafCertKeyFile,
		rootCAFiles:     rootCAFiles,
	}
}

func (t *TLS13ConfigProvider) GetTLSConfig() (*tls.Config, error) {

	leafCertificate, err := t.certProvider.LoadX509KeyPair(t.leafCertFile, t.leafCertKeyFile)
	if err != nil {
		return nil, err
	}

	rootCACertPool, err := t.certProvider.LoadX509CertPool(t.rootCAFiles)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{leafCertificate},
		RootCAs:      rootCACertPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}, nil
}
