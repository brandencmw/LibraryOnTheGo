package main

import (
	"crypto/tls"
	"libraryonthego/server/config"
	"libraryonthego/server/middleware"
	"libraryonthego/server/routes"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachAuthorRoutes(router)
	routes.AttachAuthorizationRoutes(router)

	return router
}

func setupTLS() (*tls.Config, error) {
	const rootCertFolder = "certificates"
	const serverCertFolder = "server"
	tlsConfigProvider := config.NewTLS13ConfigProvider(
		path.Join(rootCertFolder, serverCertFolder, "backend-server.crt"),
		path.Join(rootCertFolder, serverCertFolder, "backend-server.key"),
		[]string{path.Join(rootCertFolder, "root-ca.crt")},
	)

	return tlsConfigProvider.GetTLSConfig()
}

func createServer(addr string, tls *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:      ":443",
		TLSConfig: tls,
		Handler:   handler,
	}
}

func main() {

	router := setupRouter()
	tlsConfig, err := setupTLS()
	if err != nil {
		log.Fatalf("Could not configure TLS: %v", err.Error())
	}

	server := createServer(":443", tlsConfig, router)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err.Error())
	}
}
