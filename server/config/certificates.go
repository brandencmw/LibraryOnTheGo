package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path"
)

type X509CertificateProvider interface {
	LoadX509KeyPair(string, string) (tls.Certificate, error)
	LoadX509CertPool([]string) (*x509.CertPool, error)
}

type TLSCertificateProvider struct {
	certificatePath string
}

func NewTLSCertificateProvider(basePath string) *TLSCertificateProvider {
	if basePath == "" {
		basePath = "./"
	}
	return &TLSCertificateProvider{
		certificatePath: basePath,
	}
}

func (c *TLSCertificateProvider) LoadX509KeyPair(certFile, keyFile string) (tls.Certificate, error) {
	certFilePath := path.Join(c.certificatePath, certFile)
	keyFilePath := path.Join(c.certificatePath, keyFile)

	return tls.LoadX509KeyPair(certFilePath, keyFilePath)
}

func (c *TLSCertificateProvider) LoadX509CertPool(certFileNames []string) (*x509.CertPool, error) {

	certPool := x509.NewCertPool()
	for _, fileName := range certFileNames {
		certFilePath := path.Join(c.certificatePath, fileName)
		certFileBytes, err := os.ReadFile(certFilePath)
		if err != nil {
			return certPool, fmt.Errorf("Error reading %v: %v", certFilePath, err.Error())
		}

		ok := certPool.AppendCertsFromPEM(certFileBytes)
		if !ok {
			return certPool, fmt.Errorf("Failed to append %v to pool", certFilePath)
		}
	}

	return certPool, nil
}
