package main

import (
	"fmt"
	"libraryonthego/server/authentication"
	"libraryonthego/server/authors"
	"libraryonthego/server/config"
	"libraryonthego/server/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	// config.LoadEnv()
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
	router.Run(routerAddress)
}
