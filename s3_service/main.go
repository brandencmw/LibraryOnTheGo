package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"path"
	"s3/config"
	"s3/middleware"
	"s3/routes"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachAuthorRoutes(router)
	return router
}

func setupTLS() (*tls.Config, error) {
	const rootCertFolder = "certificates"
	tlsConfigProvider := config.NewMutualTLS13ConfigProvider(
		path.Join(rootCertFolder, "s3-server.crt"),
		path.Join(rootCertFolder, "s3-server.key"),
		[]string{path.Join(rootCertFolder, "root-ca.crt")},
	)
	return tlsConfigProvider.GetTLSConfig()
}

func createServer(address string, tls *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:      address,
		TLSConfig: tls,
		Handler:   handler,
	}
}

func main() {

	router := setupRouter()
	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatalf("Failed to configure TLS: %v\n", err.Error())
	}

	server := createServer(":443", tlsConfig, router)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err.Error())
	}
}
