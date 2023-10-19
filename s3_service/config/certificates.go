package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

var ServerCert tls.Certificate
var CACertPool *x509.CertPool

func LoadCertificates() {
	rootCAFile := "./certificates/root-ca.crt"
	serverCertFile := "./certificates/s3-server.crt"
	keyFile := "./certificates/s3-server.key"

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, keyFile)
	if err != nil {
		fmt.Printf("Error reading server cert: %v\n", err.Error())
	}

	fmt.Printf("Server cert private key: %T", serverCert.PrivateKey)

	rootCABytes, err := os.ReadFile(rootCAFile)
	if err != nil {
		fmt.Printf("Error reading Root CA Cert file: %v\n", err.Error())
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(rootCABytes)

	ServerCert = serverCert
	CACertPool = certPool
}
