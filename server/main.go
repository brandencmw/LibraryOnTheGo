package main

import (
	"fmt"
	"libraryonthego/server/authentication"
	"libraryonthego/server/authors"
	"libraryonthego/server/config"
	"libraryonthego/server/middleware"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func init() {
	config.DBInit()
}

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	router.GET("/authors", authors.GetAuthors)

	router.POST("/login", authentication.LoginUser)
	router.POST("/auth", middleware.AuthMiddleware, authentication.ValidateUser)

	routerAddress := fmt.Sprintf("0.0.0.0:%v", os.Getenv("SERVER_PORT"))

	certFile := path.Join("./certificates", "cert.crt")
	keyFile := path.Join("./certificates", "private.key")

	http.ListenAndServeTLS(routerAddress, certFile, keyFile, router)
}
