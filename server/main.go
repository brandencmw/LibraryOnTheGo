package main

import (
	"crypto/tls"
	"fmt"
	"libraryonthego/server/config"
	"libraryonthego/server/controllers"
	"libraryonthego/server/data"
	"libraryonthego/server/middleware"
	"libraryonthego/server/routes"
	"libraryonthego/server/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createAuthorController(client *http.Client, connPool *pgxpool.Pool) *controllers.AuthorsController {
	return controllers.NewAuthorsController(
		services.NewAuthorsService(
			data.NewS3ImageRepository(client, "https://s3_service/authors"),
			data.NewPostgresAuthorRepository(connPool),
		),
	)
}

func setupAuthorRoutes(router *gin.Engine, client *http.Client, connPool *pgxpool.Pool) {
	controller := createAuthorController(client, connPool)
	routes.AttachAuthorRoutes(router, controller)
}

func setupRouter(client *http.Client, connPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachAuthorizationRoutes(router)
	setupAuthorRoutes(router, client, connPool)
	return router
}

func setupServerTLS() (*tls.Config, error) {
	rootCACertFiles := []string{"./certificates/root-ca.crt"}
	certToKeyMap := map[string]string{
		"./certificates/server/backend-server.crt": "./certificates/server/backend-server.key",
	}
	certProvider := &config.LocalCertificateProvider{
		RootCACertFiles: rootCACertFiles,
		CertToKeyMap:    certToKeyMap,
	}

	tlsBuilder := config.TLSBuilder{CertProvider: certProvider}
	return tlsBuilder.BuildTLS(config.UseTLSVersion(tls.VersionTLS13))
}

func setupClientTLS() (*tls.Config, error) {
	rootCACertFiles := []string{"./certificates/root-ca.crt"}
	certToKeyMap := map[string]string{
		"./certificates/client/backend-client.crt": "./certificates/client/backend-client.key",
	}
	certProvider := &config.LocalCertificateProvider{
		RootCACertFiles: rootCACertFiles,
		CertToKeyMap:    certToKeyMap,
	}

	tlsBuilder := config.TLSBuilder{CertProvider: certProvider}
	return tlsBuilder.BuildTLS(config.UseTLSVersion(tls.VersionTLS13))
}

func createClient() (*http.Client, error) {
	tlsConfig, err := setupClientTLS()
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}, nil
}

func createServer(addr string, tls *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:      ":443",
		TLSConfig: tls,
		Handler:   handler,
	}
}

func main() {

	httpClient, err := createClient()
	if err != nil {
		log.Fatalf("Failed to create http client: %v", err.Error())
	}

	dbConnectionPool, err := data.CreatePostgresConnectionPool()
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err.Error())
	}
	fmt.Printf("Pool: %v", dbConnectionPool)

	router := setupRouter(httpClient, dbConnectionPool)

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
