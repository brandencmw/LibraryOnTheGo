package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path"
)

func loadServerCertificate(certFileName, keyFileName string) (tls.Certificate, error) {

	serverCertFilePath := path.Join("./certificates", certFileName)
	serverKeyFilePath := path.Join("./certificates", keyFileName)

	return tls.LoadX509KeyPair(serverCertFilePath, serverKeyFilePath)
}

func loadRootCACertPool(certFileName string) (*x509.CertPool, error) {
	rootCACertFilePath := path.Join("./certificates", certFileName)

	rootCACertBytes, err := os.ReadFile(rootCACertFilePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading Root CA Cert file: %v\n", err.Error())
	}

	rootCACertPool := x509.NewCertPool()
	ok := rootCACertPool.AppendCertsFromPEM(rootCACertBytes)

	if !ok {
		return nil, errors.New("Could not append root cert to pool")
	}

	return rootCACertPool, nil
}
