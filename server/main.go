package main

import (
	"crypto/tls"
	"libraryonthego/server/config"
	"libraryonthego/server/middleware"
	"libraryonthego/server/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachAuthorRoutes(router)
	routes.AttachAuthorizationRoutes(router)

	return router
}

func getTLSConfig() (*tls.Config, error) {
	tlsConfigProvider := config.NewTLS13ConfigProvider(
		"./certificates",
		"server/backend-server.crt",
		"server/backend-server.key",
		[]string{"root-ca.crt"},
	)

	return tlsConfigProvider.GetTLSConfig()

}

func main() {

	router := setupRouter()
	tlsConfig, err := getTLSConfig()
	if err != nil {
		log.Fatalf("Could not configure TLS: %v", err.Error())
	}

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   router,
	}

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err.Error())
	}
}
