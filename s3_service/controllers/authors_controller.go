package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadAuthorImage(c *gin.Context) {

	body, _ := c.GetRawData()
	fmt.Printf("Data: %v\n", body)

	fmt.Println("Call successful!")
	c.JSON(http.StatusOK, gin.H{
		"page": "upload author picture",
	})
}

func RetrieveAuthorImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"page": "retrieve author picture",
	})
}
