package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"s3/config"
	"s3/controllers"
	"s3/data"
	"s3/middleware"
	"s3/routes"
	"s3/services"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, client *s3.Client) {
	bucketController := controllers.NewBucketController("library-pictures", services.NewBucketService(client))
	routes.AttachBucketRoutes(router, bucketController)
}

func setupRouter(client *s3.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	setupRoutes(router, client)
	return router
}

func setupServerTLS() (*tls.Config, error) {
	rootCACertFiles := []string{"./certificates/root-ca.crt"}
	certToKeyMap := map[string]string{"./certificates/s3-server.crt": "./certificates/s3-server.key"}

	certProvider := &config.LocalCertificateProvider{
		RootCACertFiles: rootCACertFiles,
		CertToKeyMap:    certToKeyMap,
	}

	tlsConfigBuilder := config.TLSBuilder{CertProvider: certProvider}
	return tlsConfigBuilder.BuildTLS(config.UseTLSVersion(tls.VersionTLS13), config.UseMutualTLS)
}

func createServer(address string, tls *tls.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:      address,
		TLSConfig: tls,
		Handler:   handler,
	}
}

func main() {
	s3Client, err := data.CreateS3Client()
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v\n", err.Error())
	}

	router := setupRouter(s3Client)

	serverTLS, err := setupServerTLS()
	if err != nil {
		log.Fatalf("Failed to configure TLS: %v\n", err.Error())
	}

	server := createServer(":443", serverTLS, router)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err.Error())
	}
}
