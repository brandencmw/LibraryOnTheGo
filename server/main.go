package main

import (
	"libraryonthego/server/authors"
	"libraryonthego/server/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	config.DBInit()
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	router.GET("/authors", authors.GetAuthors)

	router.Run("0.0.0.0:8080")
}
