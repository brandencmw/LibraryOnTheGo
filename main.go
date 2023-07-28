package main

import (
	"LibraryOnTheGo/authors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Ping pong") })
	router.GET("/authors", authors.GetAuthors)

	router.Run("0.0.0.0:8080")
}
