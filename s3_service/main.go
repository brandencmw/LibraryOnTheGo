package main

import (
	"fmt"
	"net/http"
	"s3/config"
	"s3/controllers"
	"s3/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	config.ConfigureServerTLS()
}

func main() {

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.POST("/add-author-image", controllers.UploadAuthorImage)
	router.POST("/add-book-image")
	router.GET("/get-author-image/:object-key")
	router.GET("/get-book-image/:object-key")

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: config.ServerTLS,
		Handler:   router,
	}

	fmt.Printf("Private key: %T", server.TLSConfig.Certificates[0].PrivateKey)

	err := server.ListenAndServeTLS("", "")

	if err != nil {
		fmt.Printf("%v", err.Error())
		panic("Server failed to start")
	}
}
