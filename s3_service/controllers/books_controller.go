package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadBookImage(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"page": "upload book picture",
	})

}

func RetrieveBookImage(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"page": "retrieve book picture",
	})

}
