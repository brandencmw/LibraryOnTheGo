package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type X509CertificateProvider interface {
	LoadX509KeyPair(string, string) (tls.Certificate, error)
	LoadX509CertPool([]string) (*x509.CertPool, error)
}

type TLSCertificateProvider struct{}

func (t *TLSCertificateProvider) LoadX509KeyPair(certFilePath, keyFilePath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFilePath, keyFilePath)
}

func (t *TLSCertificateProvider) LoadX509CertPool(certFilePaths []string) (*x509.CertPool, error) {

	certPool := x509.NewCertPool()
	for _, filePath := range certFilePaths {
		certFileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return certPool, fmt.Errorf("Error reading %v: %v", filePath, err.Error())
		}

		ok := certPool.AppendCertsFromPEM(certFileBytes)
		if !ok {
			return certPool, fmt.Errorf("Failed to append %v to pool", filePath)
		}
	}

	return certPool, nil
}
