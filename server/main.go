package main

import (
	"fmt"
	"libraryonthego/server/config"
	"libraryonthego/server/middleware"
	"libraryonthego/server/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// config.DBInit()
	config.ConfigureTLS()
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.AttachAuthorRoutes(router)
	routes.AttachAuthorizationRoutes(router)

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: config.ServerTLS,
		Handler:   router,
	}

	fmt.Printf("Config: %v", server.TLSConfig)
	fmt.Printf("Private key: %T", server.TLSConfig.Certificates[0].PrivateKey)

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}
}
