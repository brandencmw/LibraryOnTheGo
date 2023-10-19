package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"libraryonthego/server/config"
	"net/http"
)

func SendAuthorImageToS3(imageBytes []byte) string {

	requestBody := []byte(fmt.Sprintf(`{
		"imageData": "%v"
	}`, imageBytes))
	fmt.Printf("Body: %v\n", requestBody)
	bodyReader := bytes.NewReader(requestBody)

	req, err := http.NewRequest(http.MethodPost, "https://s3_service/add-author-image", bodyReader)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err.Error())
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{config.ServerCert},
		RootCAs:      config.CACertPool,
		MaxVersion:   tls.VersionTLS12,
		MinVersion:   tls.VersionTLS12,
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
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error with request: %v\n", err.Error())
	}

	fmt.Printf("Status: %v\n", res.Status)
	return ""
}

func UploadAuthorInfoToDB() {

}
