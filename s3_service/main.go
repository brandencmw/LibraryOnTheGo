package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/add-author-image")
	router.POST("/add-book-image")
	router.GET("/get-author-image/:object-key")
	router.GET("/get-book-image/:object-key")

	routerAddress := fmt.Sprintf("0.0.0.0:%v", os.Getenv("S3_PORT"))
	router.Run(routerAddress)
}
