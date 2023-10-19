package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"s3/config"
	"s3/controllers"
	"s3/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadCertificates()
}

func main() {

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{config.ServerCert},
		RootCAs:      config.CACertPool,
		// ClientCAs:    config.CACertPool,
		// ClientAuth:   tls.RequireAndVerifyClientCert,
		MaxVersion: tls.VersionTLS12,
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
		GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
			for _, suite := range chi.CipherSuites {
				fmt.Println(tls.CipherSuiteName(suite))
			}

			for _, version := range chi.SupportedVersions {
				fmt.Println(tls.VersionName(version))
			}

			return nil, nil
		},
	}

	fmt.Print("Cipher suites: ")
	for _, suite := range tlsConfig.CipherSuites {
		fmt.Println(tls.CipherSuiteName(suite))
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.POST("/add-author-image", controllers.UploadAuthorImage)
	router.POST("/add-book-image")
	router.GET("/get-author-image/:object-key")
	router.GET("/get-book-image/:object-key")

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   router,
	}

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		fmt.Printf("%v", err.Error())
		panic("Server failed to start")
	}
}
