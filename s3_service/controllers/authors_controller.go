package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadAuthorImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"page": "upload author picture",
	})
}

func RetrieveAuthorImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"page": "retrieve author picture",
	})
}
