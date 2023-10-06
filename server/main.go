package main

import (
	"fmt"
	"libraryonthego/server/authentication"
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func init() {
	// config.DBInit()
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	// router.GET("/authors", controllers.GetAuthors)
	router.POST("/authors/create", middleware.AuthMiddleware, controllers.AddAuthor)

	router.POST("/login", authentication.LoginUser)
	router.POST("/auth", middleware.AuthMiddleware, authentication.ValidateUser)

	routerAddress := fmt.Sprintf("0.0.0.0:%v", os.Getenv("SERVER_PORT"))

	certFile := path.Join("./certificates", "cert.crt")
	keyFile := path.Join("./certificates", "private.key")

	err := http.ListenAndServeTLS(routerAddress, certFile, keyFile, router)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
	}
}
