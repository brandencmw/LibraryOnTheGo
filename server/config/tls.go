package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

var (
	ServerTLS  *tls.Config
	ClientTLS  *tls.Config
	rootCAPool *x509.CertPool
)

func ConfigureTLS() {
	rootCAPool, err := loadRootCACertPool("root-ca.crt")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	ServerTLS, err = configureServerTLS(rootCAPool)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	ClientTLS, err = configureClientTLS(rootCAPool)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

}

func configureServerTLS(rootCAPool *x509.CertPool) (*tls.Config, error) {

	serverCert, err := loadTLSCertificate("server/backend-server.crt", "server/backend-server.key")
	if err != nil {
		return nil, fmt.Errorf("Error loading server certificate: %s", err.Error())
	}

	serverTLS := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      rootCAPool,
		MinVersion:   tls.VersionTLS13,
	}

	return serverTLS, nil
}

func configureClientTLS(rootCAPool *x509.CertPool) (*tls.Config, error) {
	clientCert, err := loadTLSCertificate("client/backend-client.crt", "client/backend-client.key")
	if err != nil {
		return nil, fmt.Errorf("Error loading client certificate: %s", err.Error())
	}

	clientTLS := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      rootCAPool,
		MinVersion:   tls.VersionTLS13,
	}

	return clientTLS, nil
}
