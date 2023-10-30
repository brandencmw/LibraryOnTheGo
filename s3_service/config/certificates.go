package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
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

func (t *TLSCertificateProvider) LoadX509CertPool(rootCertFilePaths []string) (*x509.CertPool, error) {

	certPool := x509.NewCertPool()
	for _, filePath := range rootCertFilePaths {
		rootCACertBytes, err := os.ReadFile(filePath)
		if err != nil {
			return certPool, fmt.Errorf("Error reading Root CA Cert file: %v\n", err.Error())
		}
		ok := certPool.AppendCertsFromPEM(rootCACertBytes)
		if !ok {
			return certPool, errors.New("Could not append root cert to pool")
		}
	}

	return certPool, nil
}
