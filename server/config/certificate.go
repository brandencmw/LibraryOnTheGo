package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

var (
	ServerCert tls.Certificate
	CACertPool *x509.CertPool
)

func LoadCertificates() {
	rootCAFile := "./certificates/root-ca.crt"
	serverCertFile := "./certificates/server/backend-server.crt"
	serverKeyFile := "./certificates/server/backend-server.key"

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		fmt.Printf("Error reading server cert: %v\n", err.Error())
	}

	rootCABytes, err := os.ReadFile(rootCAFile)
	if err != nil {
		fmt.Printf("Error reading Root CA Cert file: %v\n", err.Error())
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(rootCABytes)

	ServerCert = serverCert
	CACertPool = certPool
}
