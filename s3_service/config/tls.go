package config

import (
	"crypto/tls"
	"fmt"
)

var ServerTLS *tls.Config

func ConfigureServerTLS() {

	serverCert, err := loadServerCertificate("s3-server.crt", "s3-server.key")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rootCAPool, err := loadRootCACertPool("root-ca.crt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ServerTLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      rootCAPool,
		MinVersion:   tls.VersionTLS13,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAPool,
	}

}
