package controllers

import (
	"fmt"
	"net/http"
	"s3/services"

	"github.com/gin-gonic/gin"
)

func UploadAuthorImage(c *gin.Context) {

	var request uploadAuthorImageRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Error binding request data: %v", err.Error()))
		return
	}

	err = services.UploadAuthorImageToS3Bucket(request.Headshot)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func RetrieveAuthorImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"page": "retrieve author picture",
	})
}
