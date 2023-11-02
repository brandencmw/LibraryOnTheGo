package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

type CertificateProvider interface {
	LoadServerCertificates() ([]tls.Certificate, error)
	LoadRootCertificatePool() (*x509.CertPool, error)
}

type LocalCertificateProvider struct {
	RootCACertFiles []string
	CertToKeyMap    map[string]string
}

func (c *LocalCertificateProvider) LoadServerCertificates() (certs []tls.Certificate, err error) {
	certs = make([]tls.Certificate, 0, len(c.CertToKeyMap))
	for certFile, keyFile := range c.CertToKeyMap {
		if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
			break
		} else {
			certs = append(certs, cert)
		}
	}
	return
}

func (c *LocalCertificateProvider) LoadRootCertificatePool() (certPool *x509.CertPool, err error) {
	certPool = x509.NewCertPool()
	for _, certFile := range c.RootCACertFiles {
		if certContents, err := os.ReadFile(certFile); err != nil {
			break
		} else {
			certPool.AppendCertsFromPEM(certContents)
		}
	}
	return
}
