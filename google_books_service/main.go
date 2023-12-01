package main

import (
	"crypto/tls"
	"libraryonthego/googlebooks/config"
	"libraryonthego/googlebooks/controllers"
	"libraryonthego/googlebooks/middleware"
	"libraryonthego/googlebooks/routes"
	"libraryonthego/googlebooks/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createBooksController(client *http.Client) *controllers.BooksController {
	return controllers.NewBooksController(services.NewBooksService(client))
}

func setupRouter(client *http.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachBookRotues(router, createBooksController(client))
	return router
}

func setupServerTLS() (*tls.Config, error) {
	rootCACertFiles := []string{"./certificates/root-ca.crt"}
	certToKeyMap := map[string]string{
		"./certificates/google-books-server.crt": "./certificates/google-books-server.key",
	}
	certProvider := &config.LocalCertificateProvider{
		RootCACertFiles: rootCACertFiles,
		CertToKeyMap:    certToKeyMap,
	}

	tlsBuilder := config.TLSBuilder{CertProvider: certProvider}
	return tlsBuilder.BuildTLS(config.UseTLSVersion(tls.VersionTLS13))
}

// func createClient() (*http.Client, error) {
// 	tlsConfig, err := setupClientTLS()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: tlsConfig,
// 		},
// 	}, nil
// }

func createServer(addr string, tls *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:      ":443",
		TLSConfig: tls,
		Handler:   handler,
	}
}

func main() {
	// httpClient, err := createClient()
	// if err != nil {
	// 	log.Fatalf("Failed to create http client: %v", err.Error())
	// }

	router := setupRouter(&http.Client{})

	serverTLS, err := setupServerTLS()
	if err != nil {
		log.Fatalf("Failed to initialize TLS: %v", err.Error())
	}

	server := createServer(":443", serverTLS, router)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err.Error())
	}
}
