package main

import (
	"crypto/tls"
	"fmt"
	"libraryonthego/server/authentication"
	"libraryonthego/server/config"
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// config.DBInit()
	config.LoadCertificates()
}

func main() {

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{config.ServerCert},
		RootCAs:      config.CACertPool,
	}

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	router.POST("/authors/create", middleware.AuthMiddleware, controllers.AddAuthor)

	router.POST("/login", authentication.LoginUser)
	router.POST("/auth", middleware.AuthMiddleware, authentication.ValidateUser)

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
		Handler:   router,
	}

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err.Error())
	}
}
