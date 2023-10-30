package config

import (
	"crypto/tls"
)

type TLSConfigProvider interface {
	GetTLSConfig() (*tls.Config, error)
}

type MutualTLS13ConfigProvider struct {
	certProvider    X509CertificateProvider
	leafCertFile    string
	leafCertKeyFile string
	rootCAFiles     []string
}

func NewMutualTLS13ConfigProvider(leafCertFile, leafCertKeyFile string, rootCAFiles []string) *MutualTLS13ConfigProvider {
	return &MutualTLS13ConfigProvider{
		certProvider:    &TLSCertificateProvider{},
		leafCertFile:    leafCertFile,
		leafCertKeyFile: leafCertKeyFile,
		rootCAFiles:     rootCAFiles,
	}
}

func (m *MutualTLS13ConfigProvider) GetTLSConfig() (*tls.Config, error) {

	leafCertificate, err := m.certProvider.LoadX509KeyPair(m.leafCertFile, m.leafCertKeyFile)
	if err != nil {
		return nil, err
	}

	rootCACertPool, err := m.certProvider.LoadX509CertPool(m.rootCAFiles)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{leafCertificate},
		RootCAs:      rootCACertPool,
		ClientCAs:    rootCACertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}, nil
}
