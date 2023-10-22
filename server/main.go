package main

import (
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
	fmt.Println("INITIALIZING")
	config.ConfigureTLS()
}

func main() {
	fmt.Println("HELLO THERE")
	router := gin.Default()
	fmt.Println("HELLO THERE")
	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	router.POST("/authors/create", middleware.AuthMiddleware, controllers.AddAuthor)
	router.POST("/login", authentication.LoginUser)
	router.POST("/auth", middleware.AuthMiddleware, authentication.ValidateUser)

	fmt.Printf("WHAT IS GOING ON?")

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
